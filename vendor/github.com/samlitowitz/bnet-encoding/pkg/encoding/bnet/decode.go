package bnet

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

// Unmarshal parses the BNET-encoded data and stores the result in the
// value pointed to by v. If v is nil or not a pointer, Unmarshal
// returns an InvalidUnmarshalError.
func Unmarshal(data []byte, v interface{}) error {
	var d decodeState

	d.init(data)
	return d.unmarshal(v)
}

// Unmarshaler is the interface implemented by types that can unmarshal a
// BNET description of themselves. The input can be assumed to be a valid
// encoding of a JSON value. UnmarshalBNET must copy the BNET data if it
// wishes to retain the data after returning.
type Unmarshaler interface {
	UnmarshalBNet([]byte) error
}

// An UnmarshalTypeError occurs when the attempting to unmarshal a type
// unknown to the Unmarshaller.
type UnmarshalTypeError struct {
	Value  string
	Type   reflect.Type
	Offset int64
	Struct string
	Field  string
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return "bnet: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
	}
	return "bnet: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

// An InvalidUnmarshalError occurs when attempting to unmarshal into a
// non-pointer or nil variable.
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "bnet: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "bnet: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "bnet: Unmarshal(nil " + e.Type.String() + ")"
}

func (d *decodeState) unmarshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	d.value(rv)
	return nil
}

type decodeState struct {
	data         []byte
	off          int      // read offset in data
	errorContext struct { // provides context for decode errors
		Struct string
		Field  string
	}
	structs     stringStack              // track structure stack
	savedValues map[string]reflect.Value // save values as directed by tags

	tagCache map[string]string // cache tag data
}

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.off = 0
	d.errorContext.Struct = ""
	d.errorContext.Field = ""
	d.structs = stringStack{}
	d.savedValues = make(map[string]reflect.Value)
	d.tagCache = make(map[string]string)
	return d
}

func (d *decodeState) error(err error) {
	panic(d.addErrorContext(err))
}

func (d *decodeState) addErrorContext(err error) error {
	if d.errorContext.Struct != "" || d.errorContext.Field != "" {
		switch err := err.(type) {
		case *IndexOutOfRangeError:
			err.Struct = d.errorContext.Struct
			err.Field = d.errorContext.Field
			return err
		case *InvalidSavedValueError:
			err.Struct = d.errorContext.Struct
			err.Field = d.errorContext.Field
			return err
		case *InvalidTagValueError:
			err.Struct = d.errorContext.Struct
			err.Field = d.errorContext.Field
			return err
		case *InvalidValueError:
			err.Struct = d.errorContext.Struct
			err.Field = d.errorContext.Field
			return err
		case *TagDefinitionRequiredError:
			err.Struct = d.errorContext.Struct
			err.Field = d.errorContext.Field
			return err
		case *UndefinedSavedValueError:
			err.Struct = d.errorContext.Struct
			err.Field = d.errorContext.Struct
		case *UnmarshalTypeError:
			err.Struct = d.errorContext.Struct
			err.Field = d.errorContext.Field
			return err
		}
	}
	return err
}

func (d *decodeState) value(v reflect.Value) {
	if !v.IsValid() {
		d.error(&InvalidValueError{Struct: "INVALID", Field: "INVALID"})
	}

	switch v.Kind() {
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			d.value(v.Index(i))
		}

	case reflect.Bool:
		v.SetBool(d.decodeBool())

	case reflect.Uint8:
		v.Set(reflect.ValueOf(d.decodeUint8()).Convert(v.Type()))

	case reflect.Uint16:
		v.Set(reflect.ValueOf(d.decodeUint16()).Convert(v.Type()))

	case reflect.Uint32:
		v.Set(reflect.ValueOf(d.decodeUint32()).Convert(v.Type()))

	case reflect.Uint64:
		v.Set(reflect.ValueOf(d.decodeUint64()).Convert(v.Type()))

	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		d.value(v.Elem())

	case reflect.Slice:
		d.decodeSlice(v)

	case reflect.String:
		v.SetString(d.decodeString())

	case reflect.Struct:
		d.decodeStruct(v)

	default:
		d.error(&UnmarshalTypeError{Value: "", Type: v.Type(), Offset: int64(d.off)})
	}
}

func (d *decodeState) unreadBytes(n int) {
	if d.off-n < 0 {
		d.error(&IndexOutOfRangeError{N: int64(n), Offset: int64(d.off)})
	}

	d.off = d.off - n
}

func (d *decodeState) readBytes(n int) []byte {
	if d.off+n > len(d.data) {
		d.error(&IndexOutOfRangeError{N: int64(n), Offset: int64(d.off)})
	}

	item := d.data[d.off : d.off+n]
	d.off = d.off + n

	return item
}

func (d *decodeState) readString() string {
	var n int
	for n = 0; n < len(d.data); n++ {
		if d.data[d.off+n] == 0x00 {
			break
		}
	}

	if d.off+n > len(d.data) {
		d.error(fmt.Errorf("attempting to read past buffer"))
	}

	item := d.data[d.off : d.off+n]
	d.off = d.off + n + 1 // skip null terminator

	return string(item)
}

func (d *decodeState) readStringList() []string {
	var items []string

	for item := d.readString(); len(item) != 0; item = d.readString() {
		items = append(items, item)
	}

	return items
}

func (d *decodeState) decodeBool() bool {
	var i uint32
	var unread int

	if size, ok := d.tagCache["size"]; ok {
		switch size {
		case "uint8":
			i = uint32(d.decodeUint8())
			unread = 1

		case "uint32":
			i = d.decodeUint32()
			unread = 4

		default:
			d.error(&InvalidTagValueError{Expected: "uint8 or uint32", Value: size})
		}
	} else {
		d.error(&TagDefinitionRequiredError{Tag: "size"})
	}

	if i != 0 && i != 1 {
		d.unreadBytes(unread)
		d.error(&InvalidValueError{Value: d.readBytes(unread)})
	}

	if i == 1 {
		return true
	}
	return false
}

func (d *decodeState) decodeUint8() uint8 {
	return d.readBytes(1)[0]
}

func (d *decodeState) decodeUint16() uint16 {
	if _, ok := d.tagCache["bigendian"]; ok {
		return binary.BigEndian.Uint16(d.readBytes(2))
	}
	return binary.LittleEndian.Uint16(d.readBytes(2))
}

func (d *decodeState) decodeUint32() uint32 {
	if _, ok := d.tagCache["bigendian"]; ok {
		return binary.BigEndian.Uint32(d.readBytes(4))
	}
	return binary.LittleEndian.Uint32(d.readBytes(4))
}

func (d *decodeState) decodeUint64() uint64 {
	if _, ok := d.tagCache["bigendian"]; ok {
		return binary.BigEndian.Uint64(d.readBytes(8))
	}
	return binary.LittleEndian.Uint64(d.readBytes(8))
}

func (d *decodeState) decodeSlice(v reflect.Value) {
	switch v.Type().Elem().Kind() {
	case reflect.String:
		v.Set(reflect.ValueOf(d.readStringList()))
	default:
		if len(d.tagCache) == 0 {
			d.error(&TagDefinitionRequiredError{Tag: "len"})
		}

		name, ok := d.tagCache["len"]
		if !ok {
			d.error(&TagDefinitionRequiredError{Tag: "save-" + name})
		}

		len := d.getSavedValueAsInt(name)

		sliceType := v.Type().Elem()
		if sliceType.Kind() == reflect.Ptr {
			sliceType = sliceType.Elem()
		}

		for i := 0; i < len; i++ {
			v.Set(reflect.Append(v, reflect.Zero(sliceType)))
			d.value(v.Index(i))
		}
	}
}

func (d *decodeState) decodeString() string {
	return d.readString()
}

func (d *decodeState) decodeStruct(v reflect.Value) {
	if v.Type().Name() == "" {
		// @todo AnonymousStructError
		d.error(nil)
	}
	d.errorContext.Struct = v.Type().Name()
	d.structs = d.structs.Push(v.Type().Name())

	sType := reflect.TypeOf(v.Interface())

	for i := 0; i < v.NumField(); i++ {
		d.tagCache = parseTag(sType.Field(i).Tag.Get("bnet"))
		field := v.Field(i)
		d.errorContext.Field = v.Type().Field(i).Name
		d.saveTag(field)
		d.value(field)
		d.tagCache = nil
	}

	d.structs, _ = d.structs.Pop()
}

func (d *decodeState) saveTag(v reflect.Value) {
	if name, ok := d.tagCache["save"]; ok {
		d.savedValues[name] = v
	}
}

func (d *decodeState) getSavedValueAsInt(name string) int {
	var v reflect.Value
	var ok bool
	if v, ok = d.savedValues[name]; !ok {
		d.error(&UndefinedSavedValueError{Name: name})
	}

	var i int

	switch v.Kind() {
	case reflect.Uint8:
		i = int(v.Uint())

	case reflect.Uint16:
		i = int(v.Uint())

	case reflect.Uint32:
		i = int(v.Uint())

	case reflect.Uint64:
		i = int(v.Uint())

	case reflect.String:
		t, err := strconv.Atoi(v.Interface().(string))
		if err != nil {
			d.error(&InvalidSavedValueError{Expected: "int type", Value: v.String()})
		}
		i = t

	default:
		d.error(&InvalidSavedValueError{Expected: "int type", Value: v.String()})
	}

	return i
}

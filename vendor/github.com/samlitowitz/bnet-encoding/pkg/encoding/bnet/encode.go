package bnet

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
)

// Marshal returns the BNET encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	e := &encodeState{}
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}

// Marshaler is the interface implemented by types that can marshal
// themselves into valid BNET.
type Marshaler interface {
	MarshalBNet() ([]byte, error)
}

// An UnsupportedTypeError occurs when attempting to marshal an
// unsupported type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "bnet: unsupported type: " + e.Type.String()
}

type encodeState struct {
	buf bytes.Buffer

	tagCache map[string]string
}

func (e *encodeState) Bytes() []byte {
	return e.buf.Bytes()
}

func (e *encodeState) marshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if s, ok := r.(string); ok {
				panic(s)
			}
			err = r.(error)
		}
	}()
	e.reflectValue(reflect.ValueOf(v))
	return nil
}

func (e *encodeState) error(err error) {
	panic(err)
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func (e *encodeState) reflectValue(v reflect.Value) {
	valueEncoder(v)(e, v)
}

type encoderFunc func(e *encodeState, v reflect.Value)

var encoderCache sync.Map

func valueEncoder(v reflect.Value) encoderFunc {
	if !v.IsValid() {
		return invalidValueEncoder
	}
	return typeEncoder(v.Type())
}

func typeEncoder(t reflect.Type) encoderFunc {
	if fi, ok := encoderCache.Load(t); ok {
		return fi.(encoderFunc)
	}

	var (
		wg sync.WaitGroup
		f  encoderFunc
	)
	wg.Add(1)
	fi, loaded := encoderCache.LoadOrStore(t, encoderFunc(func(e *encodeState, v reflect.Value) {
		wg.Wait()
		f(e, v)
	}))
	if loaded {
		return fi.(encoderFunc)
	}

	f = newTypeEncoder(t, true)
	wg.Done()
	encoderCache.Store(t, f)
	return f
}

var (
// marshalerType = reflect.TypeOf(new(Marshaler)).Elem()
// binaryMarshalerType = reflect.TypeOf(new(encoding.BinaryMarshaler)).Elem()
)

func newTypeEncoder(t reflect.Type, allowAddr bool) encoderFunc {
	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.String:
		return stringEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Ptr:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}

func invalidValueEncoder(e *encodeState, v reflect.Value) {
	e.error(&InvalidValueError{})
}

func boolEncoder(e *encodeState, v reflect.Value) {
	if size, ok := e.tagCache["size"]; ok {
		// @todo TagDefinitionRequiredError
		switch size {
		case "uint8":
			if v.Bool() {
				uintEncoder(e, reflect.ValueOf(uint8(0x01)))
				return
			}
			uintEncoder(e, reflect.ValueOf(uint8(0x00)))
		case "uint32":
			if v.Bool() {
				uintEncoder(e, reflect.ValueOf(uint32(0x01)))
				return
			}
			uintEncoder(e, reflect.ValueOf(uint32(0x00)))

		default:
			e.error(&InvalidTagValueError{Expected: "uint8 or uint32", Value: size})
		}
	} else {
		e.error(&TagDefinitionRequiredError{Tag: "size"})
	}
}

func uintEncoder(e *encodeState, v reflect.Value) {
	_, bigendian := e.tagCache["bigendian"]
	switch v.Kind() {
	case reflect.Uint8:
		e.buf.WriteByte(uint8(v.Uint()))
	case reflect.Uint16:
		b := []byte{0x00, 0x00}
		if bigendian {
			binary.BigEndian.PutUint16(b, uint16(v.Uint()))
		} else {
			binary.LittleEndian.PutUint16(b, uint16(v.Uint()))
		}

		e.buf.Write(b)
	case reflect.Uint32:
		b := []byte{0x00, 0x00, 0x00, 0x00}
		if bigendian {
			binary.BigEndian.PutUint32(b, uint32(v.Uint()))
		} else {
			binary.LittleEndian.PutUint32(b, uint32(v.Uint()))
		}
		e.buf.Write(b)
	case reflect.Uint64:
		b := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		if bigendian {
			binary.BigEndian.PutUint64(b, v.Uint())
		} else {
			binary.LittleEndian.PutUint64(b, v.Uint())
		}
		e.buf.Write(b)
	}
}

func stringEncoder(e *encodeState, v reflect.Value) {
	e.buf.WriteString(v.String())
	e.buf.WriteByte(0x00)
}

func interfaceEncoder(e *encodeState, v reflect.Value) {
	if v.IsNil() {
		return
	}
	e.reflectValue(v.Elem())
}

func unsupportedTypeEncoder(e *encodeState, v reflect.Value) {
	e.error(&UnsupportedTypeError{v.Type()})
}

type structEncoder struct {
	fields    []field
	fieldEncs []encoderFunc
}

func (se *structEncoder) encode(e *encodeState, v reflect.Value) {
	for i, f := range se.fields {
		fv, err := fieldByIndex(v, f.index)
		if err != nil {
			e.error(err)
		}
		if !fv.IsValid() || isEmptyValue(*fv) {
			// @todo InvalidValue
			e.error(&InvalidValueError{})
		}
		e.tagCache = se.fields[i].tags
		se.fieldEncs[i](e, *fv)
		e.tagCache = nil
	}
}

func newStructEncoder(t reflect.Type) encoderFunc {
	fields := cachedTypeFields(t)
	se := &structEncoder{
		fields:    fields,
		fieldEncs: make([]encoderFunc, len(fields)),
	}
	for i, f := range fields {
		se.fieldEncs[i] = typeEncoder(typeByIndex(t, f.index))
	}
	return se.encode
}

func encodeStringSlice(e *encodeState, v reflect.Value) {
	n := v.Len()
	for i := 0; i < n; i++ {
		stringEncoder(e, v.Index(i))
	}
	e.buf.WriteByte(0x00)
}

type sliceEncoder struct {
	arrayEnc encoderFunc
}

func (se *sliceEncoder) encode(e *encodeState, v reflect.Value) {
	if v.IsNil() {
		return
	}
	se.arrayEnc(e, v)
}

func newSliceEncoder(t reflect.Type) encoderFunc {
	if t.Elem().Kind() == reflect.String {
		return encodeStringSlice
	}
	enc := &sliceEncoder{newArrayEncoder(t)}
	return enc.encode
}

type arrayEncoder struct {
	elemEnc encoderFunc
}

func (ae *arrayEncoder) encode(e *encodeState, v reflect.Value) {
	n := v.Len()
	for i := 0; i < n; i++ {
		ae.elemEnc(e, v.Index(i))
	}
}

func newArrayEncoder(t reflect.Type) encoderFunc {
	enc := &arrayEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

func newPtrEncoder(t reflect.Type) encoderFunc {
	enc := &ptrEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

type ptrEncoder struct {
	ptrEnc encoderFunc
}

func (pe *ptrEncoder) encode(e *encodeState, v reflect.Value) {
	if v.IsNil() {
		return
	}
	pe.ptrEnc(e, v.Elem())
}

func fieldByIndex(v reflect.Value, index []int) (*reflect.Value, error) {
	for _, i := range index {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return nil, &NilPointerError{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return &v, nil
}

func typeByIndex(t reflect.Type, index []int) reflect.Type {
	for _, i := range index {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		t = t.Field(i).Type
	}
	return t
}

type field struct {
	name  string
	tags  map[string]string
	index []int
	typ   reflect.Type
}

// byIndex sorts field by index sequence.
type byIndex []field

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].index {
		if k >= len(x[j].index) {
			return false
		}
		if xik != x[j].index[k] {
			return xik < x[j].index[k]
		}
	}
	return len(x[i].index) < len(x[j].index)
}

func dominantField(fields []field) (field, bool) {
	// The fields are sorted in increasing index-length order. The winner
	// must therefore be one with the shortest index length. Drop all
	// longer entries, which is easy: just truncate the slice.
	length := len(fields[0].index)
	tagged := -1 // Index of first tagged field.
	for i, f := range fields {
		if len(f.index) > length {
			fields = fields[:i]
			break
		}
	}
	if tagged >= 0 {
		return fields[tagged], true
	}
	// All remaining fields have the same length. If there's more than one,
	// we have a conflict (two fields named "X" at the same level) and we
	// return no field.
	if len(fields) > 1 {
		return field{}, false
	}
	return fields[0], true
}

var fieldCache struct {
	value atomic.Value
	mu    sync.Mutex
}

func typeFields(t reflect.Type) []field {
	// Anonymous fields to explore at the current level and the next.
	current := []field{}
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	var count map[reflect.Type]int
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				isUnexported := sf.PkgPath != ""
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Ptr {
						t = t.Elem()
					}
					if isUnexported && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if isUnexported {
					// Ignore unexported non-embedded fields.
					continue
				}

				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					// Follow pointer.
					ft = ft.Elem()
				}

				if !sf.Anonymous || ft.Kind() != reflect.Struct {
					fields = append(fields, field{
						name:  sf.Name,
						tags:  parseTag(sf.Tag.Get("bnet")),
						index: index,
						typ:   ft,
					})
					if count[f.typ] > 1 {
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}

	sort.Slice(fields, func(i, j int) bool {
		x := fields
		// sort field by name, breaking ties with depth, then
		// breaking ties with "name came from json tag", then
		// breaking ties with index sequence.
		if x[i].name != x[j].name {
			return x[i].name < x[j].name
		}
		if len(x[i].index) != len(x[j].index) {
			return len(x[i].index) < len(x[j].index)
		}
		return byIndex(x).Less(i, j)
	})

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with JSON tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field index length. Loop over names; for each name, delete
	// hidden fields by choosing the one dominant field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.name != name {
				break
			}
		}
		if advance == 1 { // Only one field with this name
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(byIndex(fields))

	return fields
}

func cachedTypeFields(t reflect.Type) []field {
	m, _ := fieldCache.value.Load().(map[reflect.Type][]field)
	f := m[t]
	if f != nil {
		return f
	}
	f = typeFields(t)
	if f == nil {
		f = []field{}
	}

	fieldCache.mu.Lock()
	m, _ = fieldCache.value.Load().(map[reflect.Type][]field)
	newM := make(map[reflect.Type][]field, len(m)+1)
	for k, v := range m {
		newM[k] = v
	}
	newM[t] = f
	fieldCache.value.Store(newM)
	fieldCache.mu.Unlock()
	return f
}

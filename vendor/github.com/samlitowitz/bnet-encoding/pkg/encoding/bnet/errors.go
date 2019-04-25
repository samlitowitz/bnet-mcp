package bnet

import "fmt"

// An IndexOutOfRangeError describes
type IndexOutOfRangeError struct {
	N      int64
	Offset int64
	Struct string
	Field  string
}

func (e *IndexOutOfRangeError) Error() string {
	return fmt.Sprintf("bnet: index out of range when attempting to unmarshal %d bytes starting at offset %d into Go struct field %s.%s", e.N, e.Offset, e.Struct, e.Field)
}

// An InvalidSavedValueError occurs when a value is unable to be
// unmarshalled according to the respective `bnet` struct tag.
type InvalidSavedValueError struct {
	Expected string
	Value    string
	Struct   string
	Field    string
}

func (e *InvalidSavedValueError) Error() string {
	return fmt.Sprintf("bnet: invalid saved value for Go struct field %s.%s, expected %s, got %s", e.Struct, e.Field, e.Expected, e.Value)
}

// An InvalidTagValueError occurs when a `bnet` tag value is invalid.
type InvalidTagValueError struct {
	Expected string
	Value    string
	Struct   string
	Field    string
}

func (e *InvalidTagValueError) Error() string {
	return fmt.Sprintf("bnet: invalid tag value for Go struct field %s.%s, expected %s, got %s", e.Struct, e.Field, e.Expected, e.Value)
}

// An InvalidValueError occurs when a value is unable to marshalled or
// unmarshalled because it does not conform to the `bnet` tag specified.
type InvalidValueError struct {
	Struct string
	Field  string
	Value  []byte
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("bnet: invalid value for Go struct field %s.%s is %#v", e.Struct, e.Field, e.Value)
}

// A NilPointerError occurs when attempting marshal a nil pointer.
type NilPointerError struct{}

func (e *NilPointerError) Error() string {
	return "bnet: nil pointer error"
}

// A TagDefinitionRequiredError occurs when attempting to marshal or
// unmarshal a structure which is missing a required `bnet` tag.
type TagDefinitionRequiredError struct {
	Tag    string
	Struct string
	Field  string
}

func (e *TagDefinitionRequiredError) Error() string {
	return fmt.Sprintf("bnet: tag definition %s required for Go struct field %s.%s", e.Tag, e.Struct, e.Field)
}

// An UndefinedSavedValueError occurs when attempting to unmarshal into
// a struct which uses a `bnet` tag to access a saved value that is not
// defined prior.
//  type UseUndefined struct {
//      // A uint32 `bnet:"save-A"`
//      B []uint32 `bnet:"len-A"`
//  }
type UndefinedSavedValueError struct {
	Name   string
	Struct string
	Field  string
}

func (e *UndefinedSavedValueError) Error() string {
	return fmt.Sprintf("bnet: attempting to access unsaved declaration %s for Go struct field %s.%s", e.Name, e.Struct, e.Field)
}

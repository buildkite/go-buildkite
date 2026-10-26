package buildkite

import "encoding/json"

// Optional marks whether a request field should be sent in JSON.
//
// Use Optional on struct fields tagged with json:",omitzero". The zero value is
// omitted from the encoded object, and Some sends a value, including the zero
// value for T. If an Optional field is not tagged with omitzero, an unset value
// marshals as the zero value of T instead of being omitted.
type Optional[T any] struct {
	value T
	set   bool
	null  bool
}

// Some returns an Optional that sends value in JSON.
//
// Nil slices and maps are encoded as null. Use an initialized empty slice or
// map when the API expects [] or {}.
func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value: value,
		set:   true,
	}
}

// Null returns an Optional that sends JSON null. Use it only for fields the API
// accepts as nullable.
func Null[T any]() Optional[T] {
	return Optional[T]{set: true, null: true}
}

// IsZero reports whether the field should be omitted by json:",omitzero".
func (o Optional[T]) IsZero() bool {
	return !o.set
}

// Value returns the wrapped value and whether it was set.
// An explicit Null returns the zero value of T and true.
func (o Optional[T]) Value() (T, bool) {
	return o.value, o.set
}

// MarshalJSON encodes the wrapped value.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.null {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

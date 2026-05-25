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
}

// Some returns an Optional that sends value in JSON.
func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value: value,
		set:   true,
	}
}

// IsZero reports whether the field should be omitted by json:",omitzero".
func (o Optional[T]) IsZero() bool {
	return !o.set
}

// Value returns the wrapped value and whether it was set.
func (o Optional[T]) Value() (T, bool) {
	return o.value, o.set
}

// MarshalJSON encodes the wrapped value.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

// UnmarshalJSON records that the field was present and decodes its value.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.set = true
	return json.Unmarshal(data, &o.value)
}

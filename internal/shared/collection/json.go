package collection

import (
	"bytes"
	"encoding/json"

	pkg "github.com/stochastic-parrots/gollections"
)

// Marshal serializes any pkg.Collection[T] into a JSON array.
// It performs a streaming-style serialization using a bytes.Buffer to minimize
// memory allocations, avoiding the need to convert the collection to a slice first.
//
// Complexity: O(n) in time, O(n) in space for the resulting byte slice.
func Marshal[T any](c pkg.Collection[T]) ([]byte, error) {
	if c.IsEmpty() {
		return []byte("[]"), nil
	}

	var buffer bytes.Buffer
	buffer.WriteByte('[')

	first := true
	for v := range c.All() {
		if !first {
			buffer.WriteByte(',')
		}
		first = false
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		buffer.Write(b)
	}

	buffer.WriteByte(']')
	return buffer.Bytes(), nil
}

// Unmarshal populates a collection from a JSON array.
// It uses a two-step process: first decoding into a temporary slice, then
// applying the 'clear' and 'appender' functions to update the target collection.
//
// IMPORTANT: This operation is destructive. The 'clear' function is called
// before 'appender' to ensure the collection reflects exactly the JSON state.
//
// Complexity: O(n + k) where n is the current collection size and k is the JSON size.
func Unmarshal[T any](data []byte, clear func(), appender func(...T)) error {
	var temp []T
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	clear()
	appender(temp...)
	return nil
}

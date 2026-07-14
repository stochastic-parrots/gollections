package collection

import (
	"bytes"
	"encoding/json"

	"github.com/stochastic-parrots/gollections"
)

// Marshal serializes any gollections.Collection[T] into a JSON array.
// It performs a streaming-style serialization using a bytes.Buffer to minimize
// memory allocations, avoiding the need to convert the collection to a slice first.
//
// Complexity: O(n) in time, O(n) in space for the resulting byte slice.
func Marshal[T any](c gollections.Collection[T]) ([]byte, error) {
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

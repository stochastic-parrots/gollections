package collection

import (
	"fmt"
	"io"
	"iter"
	"reflect"
)

const displayLimit = 5

type formattedCollection[T any] interface {
	Enumerate() iter.Seq2[int, T]
	Length() int
}

// Format provides a standardized way to render a collection into a string.
// It is designed to be called by the Format(s fmt.State, verb rune) method of
// concrete collection implementations.
//
// Features:
//   - Supports %v, %+v (verbose), and %#v (Go-syntax).
//   - Automatically truncates output to displayLimit (5) to avoid terminal flooding.
//   - Displays logical length and physical capacity when flags are present.
//
// Complexity: O(1) in time, as it only iterates up to displayLimit elements.
func Format[T any](s fmt.State, verb rune, collection formattedCollection[T], capacity int) {
	t := reflect.TypeOf(collection)

	if verb == 'v' && s.Flag('#') {
		fmt.Fprintf(s, "%v{size:%d, cap:%d}", t, collection.Length(), capacity)
		return
	}

	if s.Flag('+') {
		fmt.Fprintf(s, "%v{len:%d, cap:%d} ", t, collection.Length(), capacity)
	}

	if collection.Length() == 0 {
		_, _ = io.WriteString(s, "[]")
		return
	}

	_, _ = io.WriteString(s, "[")
	for idx, val := range collection.Enumerate() {
		if idx >= 5 {
			fmt.Fprintf(s, " ...(+%d more)", collection.Length()-displayLimit)
			break
		}
		if idx > 0 {
			_, _ = io.WriteString(s, " ")
		}
		fmt.Fprint(s, val)
	}
	_, _ = io.WriteString(s, "]")
}

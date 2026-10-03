package collection

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/stochastic-parrots/gollections"
)

const displayLimit = 5

// Format provides a standardized way to render any gollections.Collection[T] into a string.
// It is designed to be called by the Format(s fmt.State, verb rune) method of
// concrete collection implementations.
//
// Features:
//   - Supports %v, %+v (verbose), and %#v (Go-syntax).
//   - Automatically truncates output to displayLimit (5) to avoid terminal flooding.
//   - Displays logical length and physical capacity when flags are present.
//
// Complexity: O(1) in time, as it only iterates up to displayLimit elements.
func Format[T any](s fmt.State, verb rune, collection gollections.Collection[T], capacity int) {
	t := reflect.TypeOf(collection)

	if verb == 'v' && s.Flag('#') {
		fmt.Fprintf(s, "%v{size:%d, cap:%d}", t, collection.Len(), capacity)
		return
	}

	if s.Flag('+') {
		fmt.Fprintf(s, "%v{len:%d, cap:%d} ", t, collection.Len(), capacity)
	}

	writeValues(s, collection)
}

// String renders values in iteration order using the same five-element display
// limit as Format, without type or capacity metadata.
//
// Complexity: O(1) collection operations.
func String[T any](collection gollections.Collection[T]) string {
	var output strings.Builder
	writeValues(&output, collection)
	return output.String()
}

func writeValues[T any](output io.Writer, collection gollections.Collection[T]) {
	if collection.Len() == 0 {
		_, _ = io.WriteString(output, "[]")
		return
	}

	_, _ = io.WriteString(output, "[")
	for idx, val := range collection.Enumerate() {
		if idx >= displayLimit {
			fmt.Fprintf(output, " ...(+%d more)", collection.Len()-displayLimit)
			break
		}
		if idx > 0 {
			_, _ = io.WriteString(output, " ")
		}
		fmt.Fprint(output, val)
	}
	_, _ = io.WriteString(output, "]")
}

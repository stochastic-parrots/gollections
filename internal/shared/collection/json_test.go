package collection_test

import (
	"fmt"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/shared/collection"
	"github.com/stretchr/testify/assert"
)

type ErrorMarshaler struct{}

func (e ErrorMarshaler) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("forced serialization error")
}

func TestMarshal(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		fake := &FakeCollection[int]{data: []int{}}
		got, err := collection.Marshal(fake)
		assert.NoError(t, err)
		assert.Equal(t, "[]", string(got))
	})

	t.Run("NonEmpty", func(t *testing.T) {
		fake := &FakeCollection[int]{data: []int{10, 20, 30}}
		got, err := collection.Marshal(fake)
		assert.NoError(t, err)
		assert.Equal(t, "[10,20,30]", string(got))
	})

	t.Run("Error", func(t *testing.T) {
		fake := &FakeCollection[ErrorMarshaler]{
			data: []ErrorMarshaler{{}},
		}
		_, err := collection.Marshal(fake)
		assert.Error(t, err)
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("Overwrite", func(t *testing.T) {
		currentData := []int{1, 2, 3}
		jsonData := []byte("[4,5,6]")

		clearCalled := false
		clear := func() {
			currentData = nil
			clearCalled = true
		}
		appender := func(xs ...int) {
			currentData = append(currentData, xs...)
		}

		err := collection.Unmarshal(jsonData, clear, appender)

		assert.NoError(t, err)
		assert.True(t, clearCalled, "clear should be called before appending")
		assert.Equal(t, []int{4, 5, 6}, currentData)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		err := collection.Unmarshal([]byte("[1, 2, wrong]"), func() {}, func(...int) {})
		assert.Error(t, err)
	})
}

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

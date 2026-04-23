package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDoubleLinkedNode(t *testing.T) {
	node := NewDoubleLinkedNode(5)

	assert.Equal(t, 5, node.Value)
	assert.Nil(t, node.Next)
	assert.False(t, node.HasNext())
}

func TestDoubleLinkedNode_HasNext(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	next := node.Append(6)

	assert.True(t, node.HasNext())
	assert.False(t, next.HasNext())
}

func TestDoubleLinkedNode_Append(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	next := node.Append(6)

	assert.Same(t, next, node.Next)
	assert.Same(t, node, next.Previous)
	assert.Nil(t, next.Next)
	assert.Nil(t, node.Previous)
	assert.Equal(t, 6, next.Value)
	assert.Equal(t, 5, node.Value)
}

func TestDoubleLinkedNode_Next(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	next := node.Append(6)

	assert.Same(t, next, node.Next)
	assert.Nil(t, next.Next)
}

func TestDoubleLinkedNode_HasPrevious(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	next := node.Prepend(4)

	assert.True(t, node.HasPrevious())
	assert.False(t, next.HasPrevious())
}

func TestDoubleLinkedNode_Prepend(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	previous := node.Prepend(4)
	node.Unlink()

	assert.Nil(t, node.Next)
	assert.Nil(t, previous.Previous)
	assert.Equal(t, 5, node.Value)
}

func TestDoubleLinkedNode_Unlink(t *testing.T) {
	t.Run("FullyConnected", func(t *testing.T) {
		node := NewDoubleLinkedNode(5)
		previous := node.Prepend(4)
		next := node.Append(6)
		node.Unlink()

		assert.Nil(t, node.Next)
		assert.Nil(t, node.Previous)
		assert.Equal(t, 5, node.Value)
		assert.Same(t, next, previous.Next)
		assert.Nil(t, previous.Previous)
		assert.Equal(t, 4, previous.Value)
		assert.Nil(t, next.Next)
		assert.Same(t, previous, next.Previous)
		assert.Equal(t, 6, next.Value)
	})

	t.Run("AfterPrepend", func(t *testing.T) {
		node := NewDoubleLinkedNode(5)
		previous := node.Prepend(4)
		node.Unlink()

		assert.Nil(t, node.Previous)
		assert.Nil(t, node.Next)
		assert.Equal(t, 5, node.Value)
		assert.Nil(t, previous.Previous)
		assert.Nil(t, previous.Next)
		assert.Equal(t, 4, previous.Value)
	})

	t.Run("AfterAppend", func(t *testing.T) {
		node := NewDoubleLinkedNode(5)
		next := node.Append(6)
		node.Unlink()

		assert.Nil(t, node.Previous)
		assert.Nil(t, node.Next)
		assert.Nil(t, next.Previous)
		assert.Nil(t, next.Next)
		assert.Equal(t, 5, node.Value)
		assert.Equal(t, 6, next.Value)
	})
}

func TestDoubleLinkedNode_Value(t *testing.T) {
	node := NewDoubleLinkedNode(10)
	other := NewDoubleLinkedNode("some")

	assert.Equal(t, 10, node.Value)
	assert.Equal(t, "some", other.Value)
}

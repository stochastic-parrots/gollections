package doublelinked

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDoubleLinkedNode(t *testing.T) {
	node := NewDoubleLinkedNode(5)

	assert.Equal(t, 5, node.Value())
	assert.Nil(t, node.Next())
	assert.False(t, node.HasNext())
}

func TestDoubleLinkedNodeHasNext(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	next := node.Append(6)

	assert.True(t, node.HasNext())
	assert.False(t, next.HasNext())
}

func TestDoubleLinkedNodeAppend(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	next := node.Append(6)

	assert.Same(t, next, node.Next())
	assert.Same(t, node, next.Previous())
	assert.Nil(t, next.Next())
	assert.Nil(t, node.Previous())
	assert.Equal(t, 6, next.Value())
	assert.Equal(t, 5, node.Value())
}

func TestDoubleLinkedNodeUnlinkAfterAppend(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	next := node.Append(6)
	node.Unlink()

	assert.Nil(t, node.Previous())
	assert.Nil(t, node.Next())
	assert.Nil(t, next.Previous())
	assert.Nil(t, next.Next())
	assert.Equal(t, 5, node.Value())
	assert.Equal(t, 6, next.Value())
}

func TestDoubleLinkedNodeNext(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	next := node.Append(6)

	assert.Same(t, next, node.Next())
	assert.Nil(t, next.Next())
}

func TestDoubleLinkedNodeHasPrevious(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	next := node.PreAppend(4)

	assert.True(t, node.HasPrevious())
	assert.False(t, next.HasPrevious())
}

func TestDoubleLinkedNodePreAppend(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	previous := node.PreAppend(4)
	node.Unlink()

	assert.Nil(t, node.Next())
	assert.Nil(t, previous.Previous())
	assert.Equal(t, 5, node.Value())
}

func TestDoubleLinkedNodeUnlinkAfterPreAppend(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	previous := node.PreAppend(4)
	node.Unlink()

	assert.Nil(t, node.Previous())
	assert.Nil(t, node.Next())
	assert.Equal(t, 5, node.Value())
	assert.Nil(t, previous.Previous())
	assert.Nil(t, previous.Next())
	assert.Equal(t, 4, previous.Value())
}

func TestDoubleLinkedNodeUnlink(t *testing.T) {
	node := NewDoubleLinkedNode(5)
	previous := node.PreAppend(4)
	next := node.Append(6)
	node.Unlink()

	assert.Nil(t, node.Next())
	assert.Nil(t, node.Previous())
	assert.Equal(t, 5, node.Value())
	assert.Same(t, next, previous.Next())
	assert.Nil(t, previous.Previous())
	assert.Equal(t, 4, previous.Value())
	assert.Nil(t, next.Next())
	assert.Same(t, previous, next.Previous())
	assert.Equal(t, 6, next.Value())
}

func TestDoubleLinkedNodeValue(t *testing.T) {
	node := NewDoubleLinkedNode(10)
	other := NewDoubleLinkedNode("some")

	assert.Equal(t, 10, node.Value())
	assert.Equal(t, "some", other.Value())
}

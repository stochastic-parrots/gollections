package linked

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	node := NewNode(5)

	assert.Equal(t, 5, node.Value())
	assert.Nil(t, node.Next())
	assert.False(t, node.HasNext())
}

func TestNodeHasNext(t *testing.T) {
	node := NewNode(5)
	next := node.Append(6)

	assert.True(t, node.HasNext())
	assert.False(t, next.HasNext())
}

func TestNodeAppend(t *testing.T) {
	node := NewNode(5)
	next := node.Append(6)

	assert.Same(t, next, node.Next())
	assert.Nil(t, next.Next())
	assert.Equal(t, 6, next.Value())
}

func TestNodeNext(t *testing.T) {
	node := NewNode(5)
	next := node.Append(6)

	assert.Same(t, next, node.Next())
	assert.Nil(t, next.Next())
}

func TestNodeValue(t *testing.T) {
	node := NewNode(10)
	other := NewNode("some")

	assert.Equal(t, 10, node.Value())
	assert.Equal(t, "some", other.Value())
}

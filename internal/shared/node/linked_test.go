package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLinked_Node(t *testing.T) {
	node := NewLinkedNode(5)

	assert.Equal(t, 5, node.Value)
	assert.Nil(t, node.Next)
	assert.False(t, node.HasNext())
}

func TestLinkedNode_HasNext(t *testing.T) {
	node := NewLinkedNode(5)
	next := node.Append(6)

	assert.True(t, node.HasNext())
	assert.False(t, next.HasNext())
}

func TestLinkedNode_Append(t *testing.T) {
	node := NewLinkedNode(5)
	next := node.Append(6)

	assert.Same(t, next, node.Next)
	assert.Nil(t, next.Next)
	assert.Equal(t, 6, next.Value)
}

func TestLinkedNode_Next(t *testing.T) {
	node := NewLinkedNode(5)
	next := node.Append(6)

	assert.Same(t, next, node.Next)
	assert.Nil(t, next.Next)
}

func TestLinkedNode_Value(t *testing.T) {
	node := NewLinkedNode(10)
	other := NewLinkedNode("some")

	assert.Equal(t, 10, node.Value)
	assert.Equal(t, "some", other.Value)
}

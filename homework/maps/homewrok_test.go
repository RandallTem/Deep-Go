package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	key        int
	value      int
	leftChild  *Node
	rightChild *Node
}

type OrderedMap struct {
	root *Node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	m.size++
	if m.root == nil {
		m.root = &Node{key, value, nil, nil}
		return
	}
	insert(m.root, key, value)
}

func insert(node *Node, key int, value int) {
	if key < node.key {
		if node.leftChild == nil {
			node.leftChild = &Node{key, value, nil, nil}
		} else {
			insert(node.leftChild, key, value)
		}
	} else {
		if node.rightChild == nil {
			node.rightChild = &Node{key, value, nil, nil}
		} else {
			insert(node.rightChild, key, value)
		}
	}
}

func (m *OrderedMap) Erase(key int) {
	currentNode, parentNode := findKey(m.root, key)
	switch {
	case currentNode == nil:
		return
	case currentNode.leftChild == nil && currentNode.rightChild == nil:
		if currentNode == parentNode.rightChild {
			parentNode.rightChild = nil
		} else {
			parentNode.leftChild = nil
		}
	case currentNode.leftChild == nil && currentNode.rightChild != nil:
		if currentNode == parentNode.rightChild {
			parentNode.rightChild = currentNode.rightChild
		} else {
			parentNode.leftChild = currentNode.rightChild
		}
	case currentNode.leftChild != nil && currentNode.rightChild == nil:
		if currentNode == parentNode.rightChild {
			parentNode.rightChild = currentNode.leftChild
		} else {
			parentNode.leftChild = currentNode.leftChild
		}
	case currentNode.leftChild != nil && currentNode.rightChild != nil:
		successor, successorParent := findSuccessor(currentNode)
		parentNode.rightChild = successor
		successor.leftChild = currentNode.leftChild
		successor.rightChild = currentNode.rightChild
		if successorParent != currentNode {
			successorParent.leftChild = successor.rightChild
		}
	}
	m.size--
}

func findSuccessor(node *Node) (successor *Node, parent *Node) {
	successor = node.rightChild
	for successor.leftChild != nil {
		parent = successor
		successor = successor.leftChild
	}
	return successor, parent
}

func (m *OrderedMap) Contains(key int) bool {
	searchNode, _ := findKey(m.root, key)
	return searchNode != nil
}

func findKey(currentNode *Node, key int) (node *Node, parent *Node) {
	for currentNode != nil && currentNode.key != key {
		parent = currentNode
		if currentNode.key > key {
			currentNode = currentNode.leftChild
		} else {
			currentNode = currentNode.rightChild
		}
	}
	return currentNode, parent
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	inOrderBypass(m.root, action)
}

func inOrderBypass(node *Node, action func(int, int)) {
	if node == nil {
		return
	}
	inOrderBypass(node.leftChild, action)
	action(node.key, node.value)
	inOrderBypass(node.rightChild, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

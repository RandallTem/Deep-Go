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
	if m.root == nil {
		m.root = &Node{key, value, nil, nil}
		m.size++
		return
	}
	parentNode := m.root
	currentNode := m.root
	for currentNode != nil {
		parentNode = currentNode
		if currentNode.key > key {
			currentNode = currentNode.leftChild
		} else {
			currentNode = currentNode.rightChild
		}
	}
	m.size++
	if parentNode.key > key {
		parentNode.leftChild = &Node{key, value, nil, nil}
	} else {
		parentNode.rightChild = &Node{key, value, nil, nil}
	}
}

func (m *OrderedMap) Erase(key int) {
	var parentNode *Node
	currentNode := m.root
	for currentNode != nil && currentNode.key != key {
		parentNode = currentNode
		if currentNode.key > key {
			currentNode = currentNode.leftChild
		} else {
			currentNode = currentNode.rightChild
		}
	}
	switch {
	case currentNode == nil:
		return
	case currentNode.leftChild == nil && currentNode.rightChild == nil:
		if currentNode == parentNode.rightChild {
			if parentNode != nil {
				parentNode.rightChild = nil
			} else {
				m.root = nil
			}
		} else {
			if parentNode != nil {
				parentNode.leftChild = nil
			} else {
				m.root = nil
			}
		}
	case currentNode.leftChild == nil && currentNode.rightChild != nil:
		if currentNode == parentNode.rightChild {
			if parentNode != nil {
				parentNode.rightChild = currentNode.rightChild
			} else {
				m.root = currentNode.rightChild
			}
		} else {
			if parentNode != nil {
				parentNode.leftChild = currentNode.rightChild
			} else {
				m.root = currentNode.rightChild
			}
		}
	case currentNode.leftChild != nil && currentNode.rightChild == nil:
		if currentNode == parentNode.rightChild {
			if parentNode != nil {
				parentNode.rightChild = currentNode.leftChild
			} else {
				m.root = currentNode.leftChild
			}
		} else {
			if parentNode != nil {
				parentNode.leftChild = currentNode.leftChild
			} else {
				m.root = currentNode.leftChild
			}
		}
	case currentNode.leftChild != nil && currentNode.rightChild != nil:
		successor, successorParent := findSuccessor(currentNode)
		if parentNode != nil {
			parentNode.rightChild = successor
		} else {
			m.root = successor
		}
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
	currentNode := m.root
	for currentNode != nil && currentNode.key != key {
		if currentNode.key > key {
			currentNode = currentNode.leftChild
		} else {
			currentNode = currentNode.rightChild
		}
	}
	return currentNode != nil
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

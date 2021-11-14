package dlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeNode(key string) *Node {
	return &Node{Key: key}
}

func makeRandomNode(key string) *Node {
	return makeNode("1")
}

func TestCreate(t *testing.T) {
	lst := Create()
	assert.Nil(t, lst.head)
	assert.Nil(t, lst.tail)
	assert.Zero(t, lst.size)
}

func TestInsert(t *testing.T) {
	lst := Create()

	// > - node1 - <
	node1 := makeNode("1")
	lst.Insert(node1, 0)
	assert.Equal(t, lst.head, lst.tail, node1)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node1.next)
	assert.Equal(t, lst.size, uint(1))

	// > - node2 - node1 - <
	node2 := makeNode("2")
	lst.Insert(node2, 0)
	assert.Equal(t, lst.head, node1.prev, node2)
	assert.Equal(t, node2.next, lst.tail, node1)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node1.next)
	assert.Equal(t, lst.size, uint(2))

	// > - node2 - node1 - node3 - <
	node3 := makeNode("3")
	lst.Insert(node3, lst.size)
	assert.Equal(t, lst.head, node1.prev, node2)
	assert.Equal(t, node2.next, node3.prev, node1)
	assert.Equal(t, node1.next, lst.tail, node3)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node3.next)
	assert.Equal(t, lst.size, uint(3))

	// > - node2 - node4 - node1 - node3 - <
	node4 := makeNode("4")
	lst.Insert(node4, 1)
	assert.Equal(t, lst.head, node4.prev, node2)
	assert.Equal(t, node2.next, node1.prev, node4)
	assert.Equal(t, node4.next, node3.prev, node1)
	assert.Equal(t, node1.next, lst.tail, node3)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node3.next)
	assert.Equal(t, lst.size, uint(4))

	// > - node2 - node4 - node1 - node5 - node3 - <
	node5 := makeNode("5")
	lst.Insert(node5, lst.size-1)
	assert.Equal(t, lst.head, node4.prev, node2)
	assert.Equal(t, node2.next, node1.prev, node4)
	assert.Equal(t, node4.next, node5.prev, node1)
	assert.Equal(t, node1.next, node3.prev, node5)
	assert.Equal(t, node5.next, lst.tail, node3)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node3.next)
	assert.Equal(t, lst.size, uint(5))

}

func TestPushFront(t *testing.T) {
	lst := Create()

	// > - node1 - <
	node1 := makeNode("1")
	lst.PushFront(node1)
	assert.Equal(t, lst.head, lst.tail, node1)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node1.next)
	assert.Equal(t, lst.size, uint(1))

	// > - node2 - node1 - <
	node2 := makeNode("2")
	lst.PushFront(node2)
	assert.Equal(t, lst.head, node1.prev, node2)
	assert.Equal(t, node2.next, lst.tail, node1)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node1.next)
	assert.Equal(t, lst.size, uint(2))

}

func TestPushBack(t *testing.T) {
	lst := Create()

	// > - node1 - <
	node1 := makeNode("1")
	lst.PushBack(node1)
	assert.Equal(t, lst.head, lst.tail, node1)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node1.next)
	assert.Equal(t, lst.size, uint(1))

	// > - node1 - node2 -  <
	node2 := makeNode("2")
	lst.PushBack(node2)
	assert.Equal(t, lst.head, node2.prev, node1)
	assert.Equal(t, node1.next, lst.tail, node2)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node2.next)
	assert.Equal(t, lst.size, uint(2))

}

func TestRemove(t *testing.T) {
	lst := Create()
	var err error
	var res *Node

	// > - node1 - node2 - node3 - node4 - node5 - <
	node1 := makeNode("1")
	node2 := makeNode("2")
	node3 := makeNode("3")
	node4 := makeNode("4")
	node5 := makeNode("5")
	lst.PushBack(node1)
	lst.PushBack(node2)
	lst.PushBack(node3)
	lst.PushBack(node4)
	lst.PushBack(node5)

	// > - node1 - node2 - node4 - node5 - <
	res, err = lst.Remove(node3)
	assert.Equal(t, res, node3)
	assert.Equal(t, lst.head, node2.prev, node1)
	assert.Equal(t, node1.next, node4.prev, node2)
	assert.Equal(t, node2.next, node5.prev, node4)
	assert.Equal(t, node4.next, lst.tail, node5)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node5.next)
	assert.Nil(t, node3.prev)
	assert.Nil(t, node3.next)
	assert.Nil(t, err)
	assert.Equal(t, lst.size, uint(4))

	// > - node2 - node4 - node5 - <
	res, err = lst.Remove(node1)
	assert.Equal(t, res, node1)
	assert.Equal(t, lst.head, node4.prev, node2)
	assert.Equal(t, node2.next, node5.prev, node4)
	assert.Equal(t, node4.next, lst.tail, node5)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node5.next)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node1.next)
	assert.Nil(t, err)
	assert.Equal(t, lst.size, uint(3))

	// > - node2 - node4 - <
	res, err = lst.Remove(node5)
	assert.Equal(t, res, node5)
	assert.Equal(t, lst.head, node4.prev, node2)
	assert.Equal(t, node2.next, lst.tail, node4)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node4.next)
	assert.Nil(t, node5.prev)
	assert.Nil(t, node5.next)
	assert.Nil(t, err)
	assert.Equal(t, lst.size, uint(2))

	// > - node4 - <
	res, err = lst.Remove(node2)
	assert.Equal(t, res, node2)
	assert.Equal(t, lst.head, lst.tail, node4)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node2.next)
	assert.Nil(t, node4.prev)
	assert.Nil(t, node4.next)
	assert.Nil(t, err)
	assert.Equal(t, lst.size, uint(1))

	// > - <
	res, err = lst.Remove(node4)
	assert.Equal(t, res, node4)
	assert.Equal(t, lst.head, lst.tail, nil)
	assert.Nil(t, node4.prev)
	assert.Nil(t, node4.next)
	assert.Nil(t, err)
	assert.Equal(t, lst.size, uint(0))

}

func TestPopBack(t *testing.T) {
	lst := Create()
	var err error
	var res *Node

	// > - node1 - node2 - <
	node1 := makeNode("1")
	node2 := makeNode("2")
	lst.PushBack(node1)
	lst.PushBack(node2)

	// > - node1 - <
	res, err = lst.PopBack()
	assert.Equal(t, node2, res)
	assert.Equal(t, lst.head, lst.tail, node1)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node1.next)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node2.next)
	assert.Nil(t, err)
	assert.Equal(t, lst.size, uint(1))

	// > - <
	res, err = lst.PopBack()
	assert.Equal(t, node1, res)
	assert.Equal(t, lst.head, lst.tail, nil)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node1.next)
	assert.Nil(t, err)
	assert.Equal(t, lst.size, uint(0))

}

func TestPopFront(t *testing.T) {
	lst := Create()
	var err error
	var res *Node

	// > - node1 - node2 - <
	node1 := makeNode("1")
	node2 := makeNode("2")
	lst.PushBack(node1)
	lst.PushBack(node2)

	// > - node2 - <
	res, err = lst.PopFront()
	assert.Equal(t, node1, res)
	assert.Equal(t, lst.head, lst.tail, node2)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node2.next)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node1.next)
	assert.Nil(t, err)
	assert.Equal(t, lst.size, uint(1))

	// > - <
	res, err = lst.PopFront()
	assert.Equal(t, node2, res)
	assert.Equal(t, lst.head, lst.tail, nil)
	assert.Nil(t, node2.prev)
	assert.Nil(t, node2.next)
	assert.Nil(t, err)
	assert.Equal(t, lst.size, uint(0))

}

func TestPopFromEmptyError(t *testing.T) {
	lst := Create()
	var err error

	_, err = lst.PopBack()
	assert.Equal(t, lst.head, lst.tail, nil)
	assert.Equal(t, lst.size, uint(0))
	assert.IsType(t, err, RemoveFromEmptyQueueError{})

	_, err = lst.PopFront()
	assert.Equal(t, lst.head, lst.tail, nil)
	assert.Equal(t, lst.size, uint(0))
	assert.IsType(t, err, RemoveFromEmptyQueueError{})

}

func TestNilNodeError(t *testing.T) {
	lst := Create()
	var err error

	err = lst.Insert(nil, 0)
	assert.Equal(t, lst.head, lst.tail, nil)
	assert.Equal(t, lst.size, uint(0))
	assert.IsType(t, err, NilInputNodeError{})

	err = lst.PushFront(nil)
	assert.Equal(t, lst.head, lst.tail, nil)
	assert.Equal(t, lst.size, uint(0))
	assert.IsType(t, err, NilInputNodeError{})

	err = lst.PushBack(nil)
	assert.Equal(t, lst.head, lst.tail, nil)
	assert.Equal(t, lst.size, uint(0))
	assert.IsType(t, err, NilInputNodeError{})

}

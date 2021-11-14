package dlist

import "fmt"

type Node struct {
	Key   string
	Value []byte
	prev  *Node
	next  *Node
}

type LruQueue interface {
	PushFront(*Node) error
	MoveFront(*Node) error
	Remove(*Node) (*Node, error)
	PopBack() (*Node, error)
	Size() uint
}

type List interface {
	Insert(*Node, uint) error
	PushFront(*Node) error
	PushBack(*Node) error
	Remove(*Node) (*Node, error)
	PopFront() (*Node, error)
	PopBack() (*Node, error)
	Size() uint
}

// If node from another list is passed, this may lead to very strange results
// Same problem, if single node is inserted twice
type DoubleLinkedList struct {
	head *Node
	tail *Node
	size uint
}

func CreateNode(key string, value []byte) *Node {
	return &Node{
		Key:   key,
		Value: value,
		prev:  nil,
		next:  nil,
	}
}

func Create() *DoubleLinkedList {
	return &DoubleLinkedList{
		head: nil,
		tail: nil,
		size: 0,
	}
}

func (self *DoubleLinkedList) Insert(node *Node, index uint) error {
	if node == nil {
		return NilInputNodeError{}
	}
	if index > self.size {
		return OutOfRangeError{index: index, maxIndex: self.size}
	}
	if self.size == 0 { // Empty list
		node.prev = nil
		node.next = nil
		self.head = node
		self.tail = node
	} else if index == 0 { // Push front
		node.prev = nil
		node.next = self.head
		self.head.prev = node
		self.head = node
	} else if index == self.size { // Push back
		node.prev = self.tail
		node.next = nil
		self.tail.next = node
		self.tail = node
	} else { // Insert middle
		prevNode := self.head
		for i := uint(0); i < index-1; i++ {
			prevNode = prevNode.next
		}
		nextNode := prevNode.next
		node.prev = prevNode
		node.next = nextNode
		prevNode.next = node
		nextNode.prev = node
	}
	self.size += 1
	return nil
}

func (self *DoubleLinkedList) PushFront(node *Node) error {
	return self.Insert(node, 0)
}

func (self *DoubleLinkedList) PushBack(node *Node) error {
	return self.Insert(node, self.size)
}

func (self *DoubleLinkedList) Remove(node *Node) (*Node, error) {
	if node == nil {
		return nil, NilInputNodeError{}
	}
	if self.head == nil && self.tail == nil {
		return nil, RemoveFromEmptyQueueError{}
	}

	// TODO: check backward refs, to ensure node is from this list
	prevNode := node.prev
	nextNode := node.next
	if prevNode == nil && nextNode == nil { // Only node
		self.head = nil
		self.tail = nil
		node.prev = nil
		node.next = nil
	} else if prevNode == nil { // First node
		self.head = nextNode
		nextNode.prev = nil
	} else if nextNode == nil { // Last node
		self.tail = prevNode
		prevNode.next = nil
	} else { // Middle node
		nextNode.prev = prevNode
		prevNode.next = nextNode
	}
	node.prev = nil
	node.next = nil
	self.size -= 1
	return node, nil
}

func (self *DoubleLinkedList) PopBack() (*Node, error) {
	if self.head == nil && self.tail == nil {
		return nil, RemoveFromEmptyQueueError{}
	}
	return self.Remove(self.tail)
}

func (self *DoubleLinkedList) PopFront() (*Node, error) {
	if self.head == nil && self.tail == nil {
		return nil, RemoveFromEmptyQueueError{}
	}
	return self.Remove(self.head)
}

func (self *DoubleLinkedList) MoveFront(node *Node) error {
	if node.prev == nil { // Already first
		return nil
	}
	var err error
	if _, err = self.Remove(node); err != nil {
		return err
	}
	if err = self.PushFront(node); err != nil {
		return err
	}
	return nil
}

func (self *DoubleLinkedList) Size() uint {
	return self.size
}

type OutOfRangeError struct {
	index    uint
	maxIndex uint
}

func (self OutOfRangeError) Error() string {
	return fmt.Sprintf("Index `%d` is out for range for List with size %d", self.index, self.maxIndex)
}

type RemoveFromEmptyQueueError struct{}

func (self RemoveFromEmptyQueueError) Error() string {
	return fmt.Sprintf("Queue is already empty")
}

type NilInputNodeError struct{}

func (self NilInputNodeError) Error() string {
	return fmt.Sprintf("`nil` Node is passed")
}

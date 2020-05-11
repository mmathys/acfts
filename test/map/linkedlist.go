package _map

import (
	"sync"
)

// LinkedList : Data structure
type LinkedList struct {
	mu   sync.Mutex
	Head *Node
}

// Node : A Linked List node
type Node struct {
	Next *Node
	Data *[32]byte
}

// New : Create a new Linked List
func NewList() *LinkedList {
	emptyNode := &Node{
		Next: nil,
		Data: nil,
	}
	return &LinkedList{
		Head: emptyNode,
	}
}

// If two thingies to be stored are equal.
func Equals(a *[32]byte, b *[32]byte) bool {
	for i := 0; i < 32; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Append : Appending a new node to the end of the Linked List
func (ll *LinkedList) Add(d *[32]byte) bool {
	ll.mu.Lock()
	nextNode := &Node{
		Next: nil,
		Data: d,
	}
	if ll.Head.Data == nil {
		ll.Head = nextNode
		ll.mu.Unlock()
		return true
	} else {
		mynode := ll.Head
		if Equals(mynode.Data, d) {
			ll.mu.Unlock()
			return false
		}
		for mynode.Next != nil {
			mynode = mynode.Next
			if Equals(mynode.Data, d) {
				ll.mu.Unlock()
				return false
			}
		}
		mynode.Next = nextNode
		ll.mu.Unlock()
		return true
	}
}

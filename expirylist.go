package expirylist

import "time"

type ExpiryList struct {
	latest  *Node
	oldest  *Node
	timeout time.Duration
}

func New(timeout time.Duration) *ExpiryList {
	return &ExpiryList{timeout: timeout}
}

type Node struct {
	key interface{}
	t   time.Time

	prev *Node
	next *Node
}

// NewNode creates a new node with the given key and time and adds to list in sorted order
func (el *ExpiryList) NewNode(key interface{}, t time.Time) *Node {
	node := &Node{key: key, t: t}
	el.insertSortedFromLatest(node)
	return node
}

// UpdateNode updates the time of the given node and re-inserts it in the list
func (el *ExpiryList) UpdateNode(e *Node, t time.Time) {
	e.t = t
	el.unlinkFromList(e)
	el.insertSortedFromLatest(e)
}

// DeleteNode deletes the given node from the list
func (el *ExpiryList) DeleteNode(node *Node) {
	el.unlinkFromList(node)
}

// ExpireNodes returns all nodes that have expired and removes them from the list
func (el *ExpiryList) ExpireNodes(now time.Time) (keys []interface{}) {
	for el.oldest != nil && now.Sub(el.oldest.t) >= el.timeout {
		keys = append(keys, el.oldest.key)
		el.unlinkFromList(el.oldest)
	}
	return
}

// unlinkFromList removes the given node from the list
func (el *ExpiryList) unlinkFromList(node *Node) {
	if node == nil {
		return // nothing to unlink
	}

	if node == el.latest {
		el.latest = node.prev
	}

	if node == el.oldest {
		el.oldest = node.next
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	node.next = nil
	node.prev = nil
}

// insertSortedFromLatest inserts the given node into the list in sorted order from the latest node
func (el *ExpiryList) insertSortedFromLatest(node *Node) {
	if el.latest == nil {
		el.oldest = node
		el.latest = node
		return
	}

	if node.t.After(el.latest.t) || node.t.Equal(el.latest.t) {
		node.prev = el.latest
		el.latest.next = node
		el.latest = node
		return
	}

	if node.t.Before(el.oldest.t) {
		node.next = el.oldest
		el.oldest.prev = node
		el.oldest = node
		return
	}

	// find the node after which we need to insert
	for n := el.latest; n != nil; n = n.prev {
		if node.t.After(n.t) {
			node.next = n.next
			node.prev = n
			n.next.prev = node
			n.next = node
			return
		}
	}
}

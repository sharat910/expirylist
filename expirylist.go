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

func (el *ExpiryList) NewNode(key interface{}, t time.Time) *Node {
	node := &Node{key: key, t: t, prev: el.latest}
	el.makeNodeLatest(node)
	return node
}

func (el *ExpiryList) UpdateNode(e *Node, t time.Time) {
	e.t = t
	el.getNodeToTop(e)
}

func (el *ExpiryList) DeleteNode(e *Node) {
	if e == nil {
		return // nothing to delete
	}

	if e.prev != nil {
		e.prev.next = e.next
	}

	if e.next != nil {
		e.next.prev = e.prev
	}

	if el.latest == e {
		el.latest = e.prev
	}

	if el.oldest == e {
		el.oldest = e.next
	}

	// Explicitly set to nil to ease GC
	e.next = nil
	e.prev = nil
}

func (el *ExpiryList) ExpireNodes(now time.Time) (keys []interface{}) {
	if el.oldest == nil {
		// log.Println("Map is already empty! oldest pointer nil!")
		return
	}
	for now.Sub(el.oldest.t) >= el.timeout {
		node := el.oldest
		keys = append(keys, node.key)
		el.oldest = node.next
		if node.next == nil {
			el.latest = nil
			return
		}
		node.next.prev = nil
		// Explicitly set to nil to ease GC
		node.next = nil
		node.prev = nil
	}
	return
}

func (el *ExpiryList) getNodeToTop(node *Node) {
	if el.latest == node {
		// already top most node
		return
	}

	if el.oldest != node {
		node.next.prev = node.prev
		node.prev.next = node.next

		node.prev = el.latest
		node.next = nil
		el.latest.next = node
		el.latest = node
	} else {
		el.oldest = node.next
		node.next.prev = nil

		node.prev = el.latest
		node.next = nil
		el.latest.next = node
		el.latest = node
	}
}

func (el *ExpiryList) makeNodeLatest(node *Node) {

	if el.latest == nil {
		el.oldest = node
		el.latest = node
	} else {
		// next of latest points to this
		el.latest.next = node
		// latest always points to new value
		el.latest = node
	}
}

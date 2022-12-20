package expirylist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpiryList_NewNode(t *testing.T) {
	el := New(time.Minute)
	assert.Nil(t, el.latest)
	assert.Nil(t, el.oldest)
	now := time.Now()
	node1 := el.NewNode("node1", now)
	assert.Equal(t, el.latest, node1)
	assert.Equal(t, el.oldest, node1)
	assert.Nil(t, node1.next)
	assert.Nil(t, node1.prev)
	node2 := el.NewNode("node2", now.Add(2*time.Second))
	assert.Equal(t, el.latest, node2)
	assert.Equal(t, el.oldest, node1)
	assert.Equal(t, node1.next, node2)
	assert.Equal(t, node2.prev, node1)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node2.next)
	// insert in middle
	node3 := el.NewNode("node3", now.Add(1*time.Second))
	assert.Equal(t, el.latest, node2)
	assert.Equal(t, el.oldest, node1)
	assert.Equal(t, node1.next, node3)
	assert.Equal(t, node3.prev, node1)
	assert.Equal(t, node3.next, node2)
	assert.Equal(t, node2.prev, node3)
	assert.Nil(t, node1.prev)
	assert.Nil(t, node2.next)
}

func get3NodeList(now time.Time) (el *ExpiryList, n1, n2, n3 *Node) {
	el = New(time.Minute)
	n1 = el.NewNode("node1", now.Add(time.Second))
	n2 = el.NewNode("node2", now.Add(2*time.Second))
	n3 = el.NewNode("node3", now.Add(3*time.Second))
	return
}

func TestExpiryList_DeleteNode(t *testing.T) {
	el := New(time.Minute)
	now := time.Now()

	node1 := el.NewNode("node1", now)
	el.DeleteNode(node1)
	assert.Nil(t, el.latest)
	assert.Nil(t, el.oldest)

	var n1, n2, n3 *Node
	el, n1, n2, n3 = get3NodeList(now)
	assert.Equal(t, el.latest, n3)
	assert.Equal(t, el.oldest, n1)
	el.DeleteNode(n1)
	assert.Equal(t, el.latest, n3)
	assert.Equal(t, el.oldest, n2)
	assert.Equal(t, n2.next, n3)
	assert.Equal(t, n3.prev, n2)

	el, n1, n2, n3 = get3NodeList(now)
	el.DeleteNode(n2)
	assert.Equal(t, el.latest, n3)
	assert.Equal(t, el.oldest, n1)
	assert.Equal(t, n1.next, n3)
	assert.Equal(t, n3.prev, n1)

	el, n1, n2, n3 = get3NodeList(now)
	el.DeleteNode(n3)
	assert.Equal(t, el.latest, n2)
	assert.Equal(t, el.oldest, n1)
	assert.Equal(t, n1.next, n2)
	assert.Equal(t, n2.prev, n1)
}

func TestExpiryList_UpdateNode(t *testing.T) {
	el := New(time.Minute)
	now := time.Now()

	node1 := el.NewNode("node1", now)
	el.UpdateNode(node1, now.Add(2*time.Second))
	assert.Equal(t, node1.t, now.Add(2*time.Second))
	assert.Equal(t, el.latest, node1)
	assert.Equal(t, el.oldest, node1)
	assert.Nil(t, node1.next)
	assert.Nil(t, node1.prev)

	var n1, n2, n3 *Node

	// update n3 to newer time
	el, n1, n2, n3 = get3NodeList(now)
	el.UpdateNode(n3, now.Add(10*time.Second))
	assert.Equal(t, el.latest, n3)
	assert.Equal(t, el.oldest, n1)
	assert.Equal(t, n1.next, n2)
	assert.Equal(t, n2.prev, n1)
	assert.Equal(t, n2.next, n3)
	assert.Equal(t, n3.prev, n2)

	// make n2 the latest
	el, n1, n2, n3 = get3NodeList(now)
	el.UpdateNode(n2, now.Add(10*time.Second))
	assert.Equal(t, el.latest, n2)
	assert.Equal(t, el.oldest, n1)
	assert.Equal(t, n1.next, n3)
	assert.Equal(t, n3.prev, n1)
	assert.Equal(t, n3.next, n2)
	assert.Equal(t, n2.prev, n3)
	assert.Nil(t, n1.prev)
	assert.Nil(t, n2.next)

	// make n1 the latest
	el, n1, n2, n3 = get3NodeList(now)
	el.UpdateNode(n1, now.Add(10*time.Second))
	assert.Equal(t, el.latest, n1)
	assert.Equal(t, el.oldest, n2)
	assert.Equal(t, n2.next, n3)
	assert.Equal(t, n3.prev, n2)
	assert.Equal(t, n3.next, n1)
	assert.Equal(t, n1.prev, n3)
	assert.Nil(t, n2.prev)
	assert.Nil(t, n1.next)

	// update n2 to n3 time
	el, n1, n2, n3 = get3NodeList(now)
	el.UpdateNode(n2, n2.t.Add(1*time.Second))
	assert.Equal(t, n2.t, n3.t)
	assert.Equal(t, el.latest, n2)
	assert.Equal(t, el.oldest, n1)
	assert.Equal(t, n1.next, n3)
	assert.Equal(t, n3.prev, n1)
	assert.Equal(t, n3.next, n2)
	assert.Equal(t, n2.prev, n3)
	assert.Nil(t, n1.prev)
	assert.Nil(t, n2.next)
}

func TestExpiryList_ExpireNodes(t *testing.T) {
	el := New(time.Minute)
	now := time.Now()

	el.ExpireNodes(now)
	assert.Nil(t, el.latest)
	assert.Nil(t, el.oldest)

	el.NewNode("node1", now)
	keys := el.ExpireNodes(now)
	assert.Equal(t, 0, len(keys))
	keys = el.ExpireNodes(now.Add(time.Minute))
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, "node1", keys[0])
	assert.Nil(t, el.latest)
	assert.Nil(t, el.oldest)

	var n1, n2, n3 *Node
	el, n1, n2, n3 = get3NodeList(now)
	keys = el.ExpireNodes(now.Add(time.Minute).Add(time.Second))
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, n1.key, keys[0])
	assert.Equal(t, el.latest, n3)
	assert.Equal(t, el.oldest, n2)
	assert.Equal(t, n2.next, n3)
	assert.Equal(t, n3.prev, n2)
	assert.Nil(t, n2.prev)
	assert.Nil(t, n3.next)
	keys = el.ExpireNodes(now.Add(time.Minute).Add(time.Second))
	assert.Equal(t, 0, len(keys))

	keys = el.ExpireNodes(now.Add(time.Minute).Add(3 * time.Second))
	assert.Equal(t, 2, len(keys))
	assert.Equal(t, n2.key, keys[0])
	assert.Equal(t, n3.key, keys[1])
	assert.Nil(t, el.latest)
	assert.Nil(t, el.oldest)
}

package expirylist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpiryList_NewNode(t *testing.T) {
	el := New(time.Minute)
	now := time.Now()
	node1 := el.NewNode("node1", now)
	assert.Equal(t, el.latest, node1)
	assert.Equal(t, el.oldest, node1)
	assert.Nil(t, node1.next)
	assert.Nil(t, node1.prev)
	node2 := el.NewNode("node2", now.Add(time.Second))
	assert.Equal(t, el.latest, node2)
	assert.Equal(t, el.oldest, node1)
	assert.Equal(t, node1.next, node2)
	assert.Equal(t, node2.prev, node1)
	assert.Nil(t, node1.prev)
}

func get3NodeList() (el *ExpiryList, n1, n2, n3 *Node) {
	el = New(time.Minute)
	now := time.Now()
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
	el, n1, n2, n3 = get3NodeList()
	el.DeleteNode(n1)
	assert.Equal(t, el.latest, n3)
	assert.Equal(t, el.oldest, n2)
	assert.Equal(t, n2.next, n3)
	assert.Equal(t, n3.prev, n2)

	el, n1, n2, n3 = get3NodeList()
	el.DeleteNode(n2)
	assert.Equal(t, el.latest, n3)
	assert.Equal(t, el.oldest, n1)
	assert.Equal(t, n1.next, n3)
	assert.Equal(t, n3.prev, n1)

	el, n1, n2, n3 = get3NodeList()
	el.DeleteNode(n3)
	assert.Equal(t, el.latest, n2)
	assert.Equal(t, el.oldest, n1)
	assert.Equal(t, n1.next, n2)
	assert.Equal(t, n2.prev, n1)
}

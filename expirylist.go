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

func (em *ExpiryList) NewEntry(key interface{}, t time.Time) *Node {
	entry := &Node{key: key, t: t, prev: em.latest}
	em.makeEntryLatest(entry)
	return entry
}

func (em *ExpiryList) Update(e *Node, t time.Time) {
	e.t = t
	em.getEntryToTop(e)
}

func (em *ExpiryList) ExpireEntries(now time.Time) (keys []interface{}) {
	if em.oldest == nil {
		// log.Println("Map is already empty! oldest pointer nil!")
		return
	}
	for now.Sub(em.oldest.t) >= em.timeout {
		entry := em.oldest
		keys = append(keys, entry.key)
		em.oldest = entry.next
		if entry.next == nil {
			em.latest = nil
			return
		}
		entry.next.prev = nil
	}
	return
}

func (em *ExpiryList) getEntryToTop(entry *Node) {
	if em.latest == entry {
		// already top entry
		return
	}

	if em.oldest != entry {
		entry.next.prev = entry.prev
		entry.prev.next = entry.next

		entry.prev = em.latest
		entry.next = nil
		em.latest.next = entry
		em.latest = entry
	} else {
		em.oldest = entry.next
		entry.next.prev = nil

		entry.prev = em.latest
		entry.next = nil
		em.latest.next = entry
		em.latest = entry
	}
}

func (em *ExpiryList) makeEntryLatest(entry *Node) {
	if em.oldest == nil {
		em.oldest = entry
		em.latest = entry
	} else {
		// next of latest points to this
		em.latest.next = entry
		// latest always points to new value
		em.latest = entry
	}
}

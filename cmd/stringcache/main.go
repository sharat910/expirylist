package main

import (
	"log"
	"time"

	"github.com/sharat910/expirylist"
)

type Entry struct {
	e     *expirylist.Node
	value string
}

type StringCache struct {
	em *expirylist.ExpiryList
	m  map[string]*Entry
}

func (s *StringCache) Add(key, value string, t time.Time) {
	s.m[key] = &Entry{
		e:     s.em.NewNode(key, t),
		value: value,
	}
}

func (s *StringCache) Update(key, value string, t time.Time) {
	s.m[key].value = value
	s.em.UpdateNode(s.m[key].e, t)
}

func (s *StringCache) Expire(t time.Time) {
	for _, key := range s.em.ExpireNodes(t) {
		key, ok := key.(string)
		if !ok {
			log.Fatal("type inversion failed")
		}
		log.Println("Expiring key:", key)
		delete(s.m, key)
	}
}

func NewStringCache(timeout time.Duration) *StringCache {
	return &StringCache{em: expirylist.New(timeout), m: make(map[string]*Entry)}
}

func main() {
	sc := NewStringCache(time.Minute)
	now := time.Now()
	sc.Add("hello", "hi", now)
	sc.Expire(now)
	sc.Add("ola", "hi", now)
	sc.Update("hello", "hi again", now.Add(time.Second))
	sc.Expire(now.Add(time.Minute))
}

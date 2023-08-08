package learnSync

import (
	"sync"
	"testing"
	"time"
)

type SafeMap[K comparable, V any] struct {
	data  map[K]V
	mutex sync.RWMutex
}

func (s *SafeMap[K, V]) Put(key K, val V) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = val
}

func (s *SafeMap[K, V]) Get(key K, val V) (any, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	res, ok := s.data[key]
	return res, ok

}

func (s *SafeMap[K, V]) LoadOrStore(key K, newval V) (val V, bool2 bool) {
	s.mutex.RLock()
	res, ok := s.data[key]
	s.mutex.RUnlock()
	if ok {
		return res, true
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	res, ok = s.data[key]
	if ok {
		return res, true
	}
	s.data[key] = newval
	return newval, false

}

func TestMap(t *testing.T) {
	s := &SafeMap[string, string]{
		data: make(map[string]string),
	}
	go func() {
		val1, ok := s.LoadOrStore("key1", "value1")
		t.Log("go1", val1, ok)
	}()

	go func() {
		val1, ok := s.LoadOrStore("key1", "value2")
		t.Log("go2", val1, ok)
	}()

	time.Sleep(5 * time.Second)
}

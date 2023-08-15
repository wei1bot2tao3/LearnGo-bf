package selice_queue

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

// Test
func Test(t *testing.T) {
	s := SafeResource{}
	s.AddV1(123)
	s.AddV1(234)
	fmt.Println(s.data)
}

type SafeResource struct {
	mutex sync.Mutex
	data  []any
}

func (s *SafeResource) AddV1(val any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data = append(s.data, val)
	fmt.Println(s.data)

}

func (s SafeResource) AddV2(val any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	log.Printf("%p", &s.mutex)
	s.data = append(s.data, val)
	fmt.Println(s.data)

}

type SafeResourceV2 struct {
	mutex sync.Mutex
	data  []any
}

func (s *SafeResourceV2) AddV1(val any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data = append(s.data, val)

}

func (s SafeResourceV2) AddV2(val any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data = append(s.data, val)

}

package learnSync

import (
	"fmt"
	"sync"
	"testing"
)

type MyBiz struct {
	once sync.Once
}

func (m *MyBiz) Init() {
	m.once.Do(func() {
		fmt.Println("hello word init")
	})
}

// Test
func TestMaingo(t *testing.T) {
	ac := &MyBiz{}
	ac.Init()
}

type singleton struct {
}

var s *singleton

var singletonOnce sync.Once

func GetSingleton() *singleton {
	singletonOnce.Do(func() {
		s = &singleton{}
	})
	return s
}

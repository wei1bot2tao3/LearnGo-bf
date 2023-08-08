package learnLock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Test
func TestMap(t *testing.T) {
	m := sync.Map{}
	for i := 0; i < 100; i++ {
		go func(i int) {
			m.Store(i, 1)
			for {
				v, ok := m.Load(i)
				if !ok {
					continue
				}
				m.Store(i, v.(int)+1)
				fmt.Println("i=", v)
			}
		}(i)
	}
	time.Sleep(2 * time.Second)
}

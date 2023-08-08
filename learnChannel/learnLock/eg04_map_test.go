package learnLock

import (
	"fmt"
	"testing"
	"time"
)

// Test
func TestSa(t *testing.T) {
	//fatal error: concurrent map writes
	m := map[int]int{}
	for i := 0; i < 100; i++ {
		go func() {
			for {
				v := m[i]
				m[i] = v + 1
				fmt.Println(m[i])
			}
		}()
	}
	time.Sleep(1 * time.Second)
}

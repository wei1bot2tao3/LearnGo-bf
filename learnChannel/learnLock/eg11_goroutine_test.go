package learnLock

import (
	"fmt"
	"sync"
	"testing"
)

// Test
func TestCW(t *testing.T) {
	iarr := []int{1, 2, 3, 4, 5, 6, 7}
	wg := sync.WaitGroup{}
	wg.Add(len(iarr))

	for _, iteam := range iarr {
		go func(newIteam int) {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				fmt.Println(newIteam)
			}
		}(iteam)
	}
	wg.Wait()

}

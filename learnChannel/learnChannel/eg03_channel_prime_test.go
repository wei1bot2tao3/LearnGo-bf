package learnChannel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func isPrime(num int) (isPrime bool) {
	for i := 2; i < num; i++ {
		if num%i == 0 {
			isPrime = false
			return
		}
	}
	isPrime = true
	return
}

// Test
func Test(t *testing.T) {
	startTime := time.Now()
	result := make(chan int, 200000)
	workerNumber := 10
	baseNumberChan := make(chan int, 1024)
	wg := sync.WaitGroup{}
	wg.Add(workerNumber)
	for i := 0; i < workerNumber; i++ {
		go func() {
			defer wg.Done()
			for oNum := range baseNumberChan {
				if isPrime(oNum) {
					result <- oNum
				}
			}
		}()
	}

	for num := 2; num <= 200000000; num++ {
		baseNumberChan <- num

	}
	close(baseNumberChan)
	wg.Wait()
	finishTime := time.Now()
	fmt.Println(len(result))
	fmt.Println("共耗时：", finishTime.Sub(startTime))
}

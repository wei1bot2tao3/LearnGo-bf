package learnLock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestReadBook 一个人来数5000页
func TestReadBook(t *testing.T) {
	for i := 0; i <= 10; i++ {
		NumberLock()
	}

}

func Number() {
	fmt.Println("开始数")
	var totalCount int64
	wg := sync.WaitGroup{}
	wg.Add(5000)
	for p := 0; p < 5000; p++ {
		go func() {
			defer wg.Done()
			//fmt.Print("正在统计", p, "页")
			totalCount += 100 // totalCount =totalCount+100
			// 把totalCount从盒子中取出来然后加100在放回但是没有锁就同时操作了
		}()
	}
	wg.Wait()
	fmt.Println("预计有", 100*5000, "字")
	fmt.Println(totalCount)
}

func NumberLock() {
	fmt.Println("开始数")
	var totalCount int64
	totalCountlock := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(5000)
	for p := 0; p < 5000; p++ {
		go func() {
			defer wg.Done()
			//fmt.Print("正在统计", p, "页")
			totalCountlock.Lock()
			defer totalCountlock.Unlock()
			totalCount += 100 // totalCount =totalCount+100
			// 把totalCount从盒子中取出来然后加100在放回但是没有锁就同时操作了

		}()
	}
	wg.Wait()
	fmt.Println("预计有", 100*5000, "字")
	fmt.Println(totalCount)
}

// Test
func TestLockPrice(t *testing.T) {
	fmt.Println("开始数")
	var totalCount int64
	totalCountlock := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(5000)
	for p := 0; p < 5000; p++ {
		go func() {
			defer wg.Done()
			//fmt.Print("正在统计", p, "页")
			totalCountlock.Lock()
			defer totalCountlock.Unlock()
			totalCount += 100 // totalCount =totalCount+100
			// 把totalCount从盒子中取出来然后加100在放回但是没有锁就同时操作了

		}()
	}
	wg.Wait()
	fmt.Println("预计有", 100*5000, "字")
	fmt.Println(totalCount)

}

// TestRLockPrice 读写锁
func TestRLockPrice(t *testing.T) {
	fmt.Println("开始数")
	var totalCount int64
	totalCountlock := sync.RWMutex{}
	wg := sync.WaitGroup{}
	wg.Add(5000)

	go func() {
		resuly := map[int64]struct{}{}

		for p := 0; p < 5000; p++ {
			totalCountlock.Lock()
			resuly[totalCount] = struct {
			}{}
			totalCountlock.Unlock()

		}

		fmt.Println(resuly)
	}()

	for p := 0; p < 5000; p++ {
		go func() {
			defer wg.Done()
			//fmt.Print("正在统计", p, "页")
			totalCountlock.Lock()
			defer totalCountlock.Unlock()
			totalCount += 100 // totalCount =totalCount+100
			// 把totalCount从盒子中取出来然后加100在放回但是没有锁就同时操作了

		}()

	}

	wg.Wait()
	time.Sleep(10 * time.Second)
	fmt.Println("预计有", 100*5000, "字")
	fmt.Println(totalCount)

}

// TestRLockPrice 读写锁
func TestRLockPrice2(t *testing.T) {
	fmt.Println("开始数")
	var totalCount int64
	totalCountlock := sync.RWMutex{}
	wg := sync.WaitGroup{}
	wg.Add(5000)

	resuly := map[int64]struct{}{}

	for p := 0; p < 50; p++ {
		go func(p int) {
			fmt.Println(p, "读锁开锁的时间", time.Now())
			totalCountlock.RLock()
			fmt.Println(p, "读锁拿到的时间", time.Now())
			resuly[totalCount] = struct {
			}{}
			totalCountlock.RUnlock()
		}(p)

	}

	for p := 0; p < 5; p++ {
		go func() {
			defer wg.Done()
			//fmt.Print("正在统计", p, "页")
			fmt.Println("写锁开锁的时间", time.Now())
			totalCountlock.Lock()
			fmt.Println("写锁拿到的时间", time.Now())
			defer totalCountlock.Unlock()
			totalCount += 100 // totalCount =totalCount+100
			// 把totalCount从盒子中取出来然后加100在放回但是没有锁就同时操作了

		}()

	}

	wg.Wait()
	time.Sleep(10 * time.Second)
	fmt.Println("预计有", 100*5000, "字")
	fmt.Println(totalCount)

}

package learnGronutine

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestRunPrime 计算素数
func TestRunPrime(t *testing.T) {
	startTime := time.Now()
	result := []int{}
	for num := 2; num <= 200000; num++ {
		if isPrime(num) {
			result = append(result, num)
		}
	}
	finishTime := time.Now()
	fmt.Println(len(result))
	fmt.Println("共耗时：", finishTime.Sub(startTime))
}

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

// TestRunPrimeTwoPeople 计算素数 分成两个人来计算
func TestRunPrimeTwoPeople(t *testing.T) {
	startTime := time.Now()
	result := []int{}

	go func() {
		fmt.Println("第一个开始计算", time.Now())
		result = append(result, collectPrine(2, 100000)...)
		fmt.Println("第一个完成计算", time.Now())
	}()
	go func() {
		fmt.Println("第二个开始计算", time.Now())
		result = append(result, collectPrine(100001, 200000)...)
		fmt.Println("第二个完成计算", time.Now())
	}()

	time.Sleep(5 * time.Second)
	////没有完成计算就开始打印了
	finishTime := time.Now()
	fmt.Println(len(result))
	fmt.Println("共耗时：", finishTime.Sub(startTime))
}

// collectPrine 计算素数返回 只想范围
func collectPrine(satrt int, end int) (result []int) {

	for i := satrt; i <= end; i++ {
		if isPrime(i) {
			result = append(result, i)
		}
	}
	return result
}

func TestRunPrime3(t *testing.T) {
	startTime := time.Now()
	result := []int{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Println("第一个worker开始计算", time.Now())
		result = append(result, collectPrine(2, 100000)...)
		fmt.Println("第一个worker完成计算", time.Now())
	}()
	go func() {
		defer wg.Done()
		fmt.Println("第二个worker开始计算", time.Now())
		result = append(result, collectPrine(100001, 200000)...)
		fmt.Println("第二个worker完成计算", time.Now())
	}()
	wg.Wait()
	finishTime := time.Now()
	fmt.Println("finishTime: ", finishTime)
	fmt.Println(len(result))
	fmt.Println("共耗时：", finishTime.Sub(startTime))
}

// TestHelloGoroutine
func TestHelloGoroutine(t *testing.T) {
	go fmt.Println("hellw goroutine")
	time.Sleep(1 * time.Second)
}

func TestHelloGoroutine2(t *testing.T) {
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println(i)
			time.Sleep(1 * time.Second)
		}
	}()

	for i := 100; i <= 110; i++ {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}

}

// TestTestFormergoroutine
func TestFormergoroutine(t *testing.T) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			go func() {
				fmt.Println("启动新的routine@", time.Now())
				time.Sleep(1 * time.Hour)
			}()
		}
	}()

	for {
		//函数用于获取当前程序中正在运行的goroutine的数量
		fmt.Println(runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
	}

}

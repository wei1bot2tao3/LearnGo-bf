package learnChannel

import (
	"fmt"
	"testing"
	"time"
)

// Test
func TestTChannel(t *testing.T) {

	var TChannel chan int
	var TMap map[string]int = map[string]int{}

	fmt.Println(TMap)
	//都是空同时都不能写
	TChannel = make(chan int, 1)
	TChannel <- 3
	out := TChannel
	out2, ok := <-TChannel
	fmt.Println(ok)
	fmt.Println(out, out2)

}

func TestChanPutGet(t *testing.T) {
	intCh := make(chan int) //创建一个不带size的channel 不带buffer
	woekerCount := 10
	for i := 0; i < woekerCount; i++ {
		go func(i int) {
			intCh <- i
		}(i)
	}

	for o := 0; o < woekerCount; o++ {
		go func(o int) {
			out := <-intCh
			fmt.Println("出,", o, "拿到", out)
		}(o)
	}
	time.Sleep(1 * time.Second)

}

func TestChanPutGet_0wait(t *testing.T) {
	intCh := make(chan int) //创建一个不带size的channel 不带buffer
	woekerCount := 10
	for i := 0; i < woekerCount; i++ {
		go func(i int) {
			fmt.Println("开始工作", time.Now())
			intCh <- i
			fmt.Println("结束工作", time.Now())
		}(i)
	}
	time.Sleep(5 * time.Second)

	for o := 0; o < woekerCount; o++ {
		go func(o int) {
			out := <-intCh
			fmt.Println("出,", o, time.Now(), "拿到", out)
		}(o)
	}
	time.Sleep(1 * time.Second)
}

func TestChanPutGet_10wait(t *testing.T) {
	intCh := make(chan int, 10) //创建一个不带size的channel 不带buffer
	woekerCount := 10
	for i := 0; i < woekerCount; i++ {
		go func(i int) {
			fmt.Println("开始工作", time.Now())
			intCh <- i
			fmt.Println("结束工作", time.Now())
		}(i)
	}
	time.Sleep(5 * time.Second)

	for o := 0; o < woekerCount; o++ {
		go func(o int) {
			out := <-intCh
			fmt.Println("出,", o, time.Now(), "拿到", out)
		}(o)
	}
	time.Sleep(1 * time.Second)
}

func TestRangeChannel(t *testing.T) {
	intCh := make(chan int, 10)
	intCh <- 1
	intCh <- 2
	intCh <- 3
	intCh <- 4
	intCh <- 5
	intCh <- 6
	intCh <- 7
	intCh <- 8
	o1, ok := <-intCh
	// 假数据
	if ok {
		fmt.Println(o1)
	}
	fmt.Println(ok)
	close(intCh)

	for v := range intCh {
		fmt.Println(v)
	}

}

// Test
func TestSelectChannel(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan string)
	ch1 <- 1
	//go func() {
	//	time.Sleep(1 * time.Second)
	//	ch1 <- 1
	//}()
	//time.Sleep(1 * time.Second)
	go func() {
		ch2 <- "GoGOGO"
	}()

	select {
	case o := <-ch1:
		fmt.Println("ch1done", o)
	case o := <-ch2:
		fmt.Println("ch2done", o)
	default:
		fmt.Println("no done")
	}

	fmt.Println("done")
}

// Test
func TestSelectDefaultChannelandCloss(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan string)

	//go func() {
	//	time.Sleep(1 * time.Second)
	//	ch1 <- 1
	//}()
	//time.Sleep(1 * time.Second)
	close(ch1)
	go func() {
		ch2 <- "GoGOGO"
	}()

	select {
	case o := <-ch1:
		fmt.Println("ch1done", o)
	case o := <-ch2:
		fmt.Println("ch2done", o)
	default:
		fmt.Println("no done")
	}

	fmt.Println("done")
}

// Test
func TestMultipleChannel(t *testing.T) {
	ch1 := make(chan int, 1)

	for i := 0; i < 10; i++ {
		go func(i int) {
			select {
			case <-ch1:
				fmt.Println(time.Now(), i)
			}
		}(i)
	}

	fmt.Println("关channel", time.Now())
	close(ch1)
	time.Sleep(1 * time.Second)
}

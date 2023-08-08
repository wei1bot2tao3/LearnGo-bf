package learnChannel

import (
	"fmt"
	"testing"
)

// Test
func TestWithDirection(t *testing.T) {
	c := make(chan int, 100)
	inOnly(c)
	outOnly(c)
}

func inOnly(c chan<- int) {
	c <- 1
	//<-c //单项入只能是进不能取会编译错误
}
func outOnly(c <-chan int) {
	//c<-1// 相反
	<-c
}

// TestEg03 注意事项测试
func TestEg03(t *testing.T) {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	close(ch2)
	close(ch1)
	ch1Counter, ch2Counter := 0, 0
	for i := 0; i <= 100000000; i++ {
		select {
		case <-ch1:
			ch1Counter++
		case <-ch2:
			ch2Counter++

		}
	}
	fmt.Println(ch2Counter, ch1Counter)

}

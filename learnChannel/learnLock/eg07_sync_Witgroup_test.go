package learnLock

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// TestWaitGroup
func TestWaitGroup(t *testing.T) {
	runnerCount := 10
	runners := []Runner{}
	startPointWg := sync.WaitGroup{}
	startPointWg.Add(1)
	rwg := sync.WaitGroup{}
	rwg.Add(10)

	for i := 0; i < runnerCount; i++ {
		runners = append(runners, Runner{
			Name: fmt.Sprintf("%d", i),
		})
	}
	for _, items := range runners {

		go func(r Runner) {
			r.Run(&startPointWg, &rwg)
		}(items)

		//这是因为在循环中创建goroutine时，使用了闭包引用了循环变量items，而不是创建一个新的变量。
		//在循环的每次迭代中，都会创建一个新的goroutine，并且这些goroutine共享同一个items变量。当goroutine开始执行时，它们引用的items变量的值已经是循环结束时的最后一个值，即9。因此，所有的goroutine都打印出了9作为选手的名字。
	}
	startPointWg.Done()
	rwg.Wait()

}

type Runner struct {
	Name string
}

func (r Runner) Run(startPointWg, rwg *sync.WaitGroup) {
	defer rwg.Done()
	startPointWg.Wait()
	satrt := time.Now()

	fmt.Println(r.Name, "go", satrt)
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Uint64()%10) * time.Second)
	finsh := time.Now()
	fmt.Println(r.Name, "到终点了", finsh.Sub(satrt))
}

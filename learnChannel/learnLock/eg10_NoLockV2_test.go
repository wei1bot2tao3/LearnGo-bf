package learnLock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type WebServer2 struct {
	config *Config2
}

type Config2 struct {
	Content string
}

func (ws *WebServer2) Visit2() string {
	return ws.config.Content
}
func (ws *WebServer2) Reload2() {
	ws.config = &Config2{
		Content: fmt.Sprintf("%d", time.Now().UnixNano()),
	}

}

func (ws *WebServer2) ReloadWorker() {
	for {
		//time.Sleep(10 * time.Millisecond)
		ws.Reload2()
	}
}

// Test
func TestNoLock2(t *testing.T) {
	v2 := &WebServer2{
		config: &Config2{
			Content: "a",
		},
	}
	go v2.ReloadWorker()
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				v2.Visit2()
			}
		}()
	}
	wg.Wait()
	fish := time.Now()
	fmt.Println("消耗总时间为", fish.Sub(start))

}

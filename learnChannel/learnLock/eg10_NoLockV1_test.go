package learnLock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type WebServer1 struct {
	config     Config
	configLook sync.RWMutex
}

type Config struct {
	Content string
}

func (ws *WebServer1) Visit() string {
	ws.configLook.RLock()
	defer ws.configLook.RUnlock()
	return ws.config.Content
}
func (ws *WebServer1) Reload() {
	ws.configLook.Lock()
	defer ws.configLook.Unlock()
	ws.config.Content = fmt.Sprintf("%d", time.Now().UnixNano())
}

func (ws *WebServer1) ReloadWorker() {
	for {
		time.Sleep(10 * time.Millisecond)
		ws.Reload()
	}
}

// Test
func TestNoLock(t *testing.T) {
	v1 := &WebServer1{}
	go v1.ReloadWorker()
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				v1.Visit()
			}
		}()
	}
	wg.Wait()
	fish := time.Now()
	fmt.Println("消耗总时间为", fish.Sub(start))

}

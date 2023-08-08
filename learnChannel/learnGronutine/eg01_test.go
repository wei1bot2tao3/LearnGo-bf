package learnGronutine

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"
)

// TestReadBook 一个人来数5000页
func TestReadBook(t *testing.T) {
	fmt.Println("开始数")
	var totalCount int64
	for p := 0; p <= 5000; p++ {
		fmt.Println("正在统计", p, "页")
		time.Sleep(1 * time.Second)
		r, _ := rand.Int(rand.Reader, big.NewInt(800))
		fmt.Println("这一页有", r)
		totalCount += r.Int64()
	}
	fmt.Println(totalCount)

}

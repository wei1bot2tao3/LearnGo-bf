package learnLock

import (
	"sync"
	"testing"
)

type Rank struct {
	standard []string
}

var golalRank = &Rank{}
var once sync.Once = sync.Once{}

// Test
func initGobalRankStandard(standard []string) {
	once.Do(func() {
		golalRank.standard = standard
	})

}

// Test
func Test(t *testing.T) {
	sandard := []string{"asin"}
	for i := 0; i < 10; i++ {
		go func() {
			initGobalRankStandard(sandard)
		}()
	}
}

type Conn struct {
}

// DBFactory 数据库工厂
type DBFactory interface {
	GetConnection() *Conn
}

var facStore = &dbFactoryStore{}

type dbFactoryStore struct {
	store map[string]DBFactory
}

// 初始化Mysql
func initMysql(connStr string) DBFactory {

	return &MySqlDBFactory{}
}

type MySqlDBFactory struct {
	once sync.Once
}

func (MySqlDBFactory) GetConnection() *Conn {
	once.Do(func() {
		initMysql("")
	})

	return nil
}

func TestMYSql(t *testing.T) {

}
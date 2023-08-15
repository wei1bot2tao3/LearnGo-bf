package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

// 典型的 Option 设计模式
type Option func(*App)

// ShutdownCallback 采用 context.Context 来控制超时，而不是用 time.After 是因为
// - 超时本质上是使用这个回调的人控制的
// - 我们还希望用户知道，他的回调必须要在一定时间内处理完毕，而且他必须显式处理超时错误
type ShutdownCallback func(ctx context.Context)

// 你需要实现这个方法
func WithShutdownCallbacks(cbs ...ShutdownCallback) Option {

	return func(app *App) {
		app.cbs = cbs
	}
	//panic("implement me")
}

// 这里我已经预先定义好了各种可配置字段
type App struct {
	servers []*Server

	// 优雅退出整个超时时间，默认30秒
	shutdownTimeout time.Duration

	// 优雅退出时候等待处理已有请求时间，默认10秒钟
	waitTime time.Duration
	// 自定义回调超时时间，默认三秒钟
	cbTimeout time.Duration

	cbs []ShutdownCallback
}

// NewApp 创建 App 实例，注意设置默认值，同时使用这些选项
func NewApp(servers []*Server, opts ...Option) *App {
	return &App{
		servers:         servers,
		shutdownTimeout: 30 * time.Second,
		waitTime:        10 * time.Second,
		cbTimeout:       3 * time.Second,
	}

}

// StartAndServe 你主要要实现这个方法
func (app *App) StartAndServe() {
	for _, s := range app.servers {
		srv := s
		go func() {
			// 这里就启动了
			if err := srv.Start(); err != nil {
				if err == http.ErrServerClosed {
					log.Printf("服务器%s已关闭", srv.name)
				} else {
					log.Printf("服务器%s异常退出", srv.name)
				}
			}

		}()
	}
	// 从这里开始优雅退出监听系统信号，强制退出以及超时强制退出。
	wg := sync.WaitGroup{}
	wg.Add(1)
	c := make(chan os.Signal, 2)
	// signal.Notify ：c：一个通道，用于接收信号通知。 sig：一个可变参数，表示要监听的信号列表。
	//这里是把监控到的信号发送到c里面 缓冲区是1
	signal.Notify(c, signals...)
	// 开一个goroutine来监控信号
	// 然后还要在开一goroutine吗？ 监控退出，如过超时会怎么版 所以开俩go 那么缓冲区就得变成2
	//再来一次就强制退出
	go func() {
		select {
		// 先阻塞住 ，当有人请求关闭才能先开第二个go 监控然后在进行优雅
		case <-c:
			// 拒绝新请求放到 关闭函数里
			fmt.Println("开始执行优雅退出")
			go func() {
				select {
				case <-c:
					fmt.Println("强制结束")
					//wg.Done()
					os.Exit(1)
				case <-time.After(app.shutdownTimeout):
					fmt.Println("超时强制结束")
					//wg.Done()
					os.Exit(1)
				}
			}()
			app.shutdown()
			wg.Done()
		}
		//单独开一个go来监控强制退出或者超时
	}()
	wg.Wait()

	// 优雅退出的具体步骤在 shutdown 里面实现
	// 所以你需要在这里恰当的位置，调用 shutdown
	// 还剩具体实现shutdown里面的优雅了
}

// shutdown 你要设计这里面的执行步骤。
func (app *App) shutdown() {
	log.Println("开始关闭应用，停止接收新请求")
	// 你需要在这里让所有的 server 拒绝新请求
	for _, ser := range app.servers {
		ser.rejectReq()
		fmt.Println("停止接收新请求成功", ser.mux.reject)
	}
	log.Println("等待正在执行请求完结")
	// 在这里等待一段时间

	log.Println("开始关闭服务器")

	// 并发关闭服务器，同时要注意协调所有的 server 都关闭之后才能步入下一个阶段
	wg := sync.WaitGroup{}
	wg.Add(cap(app.servers))
	for _, ser := range app.servers {
		go func(ser *Server) {
			err := ser.stop(context.Background())
			if err != nil {
				wg.Done()
			}
		}(ser)
	}
	wg.Wait()
	log.Println("开始执行自定义回调")
	// 并发执行回调，要注意协调所有的回调都执行完才会步入下一个阶段
	wgcb := sync.WaitGroup{}
	wgcb.Add(cap(app.cbs))
	for _, cb := range app.cbs {
		go func(cb ShutdownCallback) {
			//context.WithTimeout() 创建一个超时的上下文 ，返回(Context, CancelFunc)
			//并返回一个派生的上下文（derived context）和一个取消函数（cancel function）
			//app.cbTimeout 目前是固定的 不然要在go签名声明一下不然会导致取到错误时间
			ctx, cschan := context.WithTimeout(context.Background(), app.cbTimeout)
			cb(ctx)
			defer cschan()
			wgcb.Done()
		}(cb)
	}

	// 释放资源
	log.Println("开始释放资源")
	wg.Wait()
	// 这一个步骤不需要你干什么，这是假装我们整个应用自己要释放一些资源
	app.close()

}

func (app *App) close() {
	// 在这里释放掉一些可能的资源
	time.Sleep(time.Second)
	log.Println("应用关闭")
}

// Server 本身可以是很多种 Server，例如 http server
// 或者 RPC server
// 理论上来说，如果你设计一个脚手架的框架，那么 Server 应该是一个接口
type Server struct {
	srv  *http.Server
	name string
	mux  *serverMux
}

// serverMux 既可以看做是装饰器模式，也可以看做委托模式
type serverMux struct {
	reject bool
	*http.ServeMux
}

// 判断一下如果失败会返回啥
func (s *serverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 只是在考虑到 CPU 高速缓存的时候，会存在短时间的不一致性
	if s.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务已关闭"))
		return
	}
	s.ServeMux.ServeHTTP(w, r)
}

func NewServer(name string, addr string) *Server {
	mux := &serverMux{ServeMux: http.NewServeMux()}
	return &Server{
		name: name,
		mux:  mux,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// Start 开始监听这 个地址
func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

// rejectReq() 声明时候默认是false 所以该成true就好了
func (s *Server) rejectReq() {
	s.mux.reject = true
}

//func (s *Server) rejectClose() {
//	s.mux.reject = false
//}

func (s *Server) stop(ctx context.Context) error {
	log.Printf("服务器%s关闭中", s.name)
	// 在这里模拟停下服务器 http包的 server.Shutdown(ctx)
	return s.srv.Shutdown(ctx)
}

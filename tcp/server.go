package tcp

import (
	"context"
	"fmt"
	"godis/interface/tcp"
	"godis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/*
@author: shg
@since: 2023/2/23 11:43 PM
@mail: shgang97@163.com
*/

// Config tcp server 属性配置
type Config struct {
	Address    string        `yaml:"address"`
	MaxConnect uint32        `yaml:"max-connect"`
	Timeout    time.Duration `yaml:"timeout"`
}

// ListenAndServeWithSignal 监听中断信号并通过 closeChan 通知服务器关闭
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		logger.Info(sig)
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Info(fmt.Sprintf("binding: %s, start listening...", cfg.Address))
	ListenAndServe(listener, handler, closeChan)
	return nil
}

// ListenAndServe 监听并提供服务，并在收到 closeChan 发来的关闭通知后关闭
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	errCh := make(chan error)
	defer close(errCh)

	// 启动一个 goroutine，当收到 closeChan 或 errChan 时进行关闭
	go func() {
		select {
		case <-closeChan:
			logger.Info("get exit signal")
			//case err := <-errCh:
			//	logger.Info(fmt.Sprintf("accept error: %s", err.Error()))
		}
		logger.Info("shutting down...")
		_ = listener.Close() // 停止监听，listener.Accept()会立即返回 io.EOF
		_ = handler.Close()  // 关闭连接
	}()

	ctx := context.Background()
	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			//errCh <- err
			break
		}

		logger.Info("accept link")
		wg.Add(1)
		// 启动一个 goroutine 来处理连接请求
		go func() {
			defer wg.Done()
			handler.Handle(ctx, conn)
		}()
	}
	wg.Wait()
}

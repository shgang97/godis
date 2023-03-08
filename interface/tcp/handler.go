package tcp

import (
	"context"
	"net"
)

/*
@author: shg
@since: 2023/2/23 11:44 PM
@mail: shgang97@163.com
*/

// Handler 是应用层服务器的抽象
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

// HandleFunc 是 handle 函数的抽象
type HandleFunc func(ctx context.Context, conn net.Conn)

package handler

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
)

/*
@author: shg
@since: 2023/3/9 3:42 AM
@mail: shgang97@163.com
*/

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type GodisHandler struct {
	activeConn sync.Map
	//db database.DB
	closing atomic.Bool
}

func (h *GodisHandler) Handle(ctx context.Context, conn net.Conn) {

}

func (h *GodisHandler) Close() error {
	return nil
}

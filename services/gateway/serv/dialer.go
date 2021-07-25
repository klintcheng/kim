package serv

import (
	"net"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/tcp"
	"github.com/klintcheng/kim/wire/pkt"
	"google.golang.org/protobuf/proto"
)

type TcpDialer struct {
	ServiceId string
}

func NewDialer(serviceId string) kim.Dialer {
	return &TcpDialer{
		ServiceId: serviceId,
	}
}

// DialAndHandshake(context.Context, string) (net.Conn, error)
func (d *TcpDialer) DialAndHandshake(ctx kim.DialerContext) (net.Conn, error) {
	// 1. 拨号建立连接
	conn, err := net.DialTimeout("tcp", ctx.Address, ctx.Timeout)
	if err != nil {
		return nil, err
	}
	req := &pkt.InnerHandshakeReq{
		ServiceId: d.ServiceId,
	}
	logger.Infof("send req %v", req)
	// 2. 把自己的ServiceId发送给对方
	bts, _ := proto.Marshal(req)
	err = tcp.WriteFrame(conn, kim.OpBinary, bts)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

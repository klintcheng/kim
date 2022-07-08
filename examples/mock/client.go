package mock

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/tcp"
	"github.com/klintcheng/kim/websocket"
)

// ClientDemo Client demo
type ClientDemo struct {
}

func (c *ClientDemo) Start(userID, protocol, addr string) {
	var cli kim.Client

	// step1: 初始化客户端
	if protocol == "ws" {
		cli = websocket.NewClient(userID, "client", websocket.ClientOptions{})
		// set dialer
		cli.SetDialer(&WebsocketDialer{})
	} else if protocol == "tcp" {
		cli = tcp.NewClient("test1", "client", tcp.ClientOptions{})
		cli.SetDialer(&TCPDialer{})
	}

	// step2: 建立连接
	err := cli.Connect(addr)
	if err != nil {
		logger.Error(err)
		return
	}
	count := 10
	go func() {
		// step3: 发送消息然后退出
		for i := 0; i < count; i++ {
			err := cli.Send([]byte(fmt.Sprintf("hello_%d", i)))
			if err != nil {
				logger.Error(err)
				return
			}
			time.Sleep(time.Millisecond * 10)
		}
	}()

	// step4: 接收消息
	recv := 0
	for {
		frame, err := cli.Read()
		if err != nil {
			logger.Info(err)
			break
		}
		if frame.GetOpCode() != kim.OpBinary {
			continue
		}
		recv++
		logger.Infof("%s receive message [%s]", cli.ServiceID(), frame.GetPayload())
		if recv == count { // 接收完消息
			break
		}
	}
	//退出
	cli.Close()
}

// WebsocketDialer WebsocketDialer
type WebsocketDialer struct {
}

// DialAndHandshake DialAndHandshake
func (d *WebsocketDialer) DialAndHandshake(ctx kim.DialerContext) (net.Conn, error) {
	logger.Info("start ws dial: ", ctx.Address)
	// 1 调用ws.Dial拨号
	ctxWithTimeout, cancel := context.WithTimeout(context.TODO(), ctx.Timeout)
	defer cancel()

	conn, _, _, err := ws.Dial(ctxWithTimeout, ctx.Address)
	if err != nil {
		return nil, err
	}
	// 2. 发送用户认证信息，示例就是userid
	err = wsutil.WriteClientBinary(conn, []byte(ctx.Id))
	if err != nil {
		return nil, err
	}
	// 3. return conn
	return conn, nil
}

// TCPDialer TCPDialer
type TCPDialer struct {
}

// DialAndHandshake DialAndHandshake
func (d *TCPDialer) DialAndHandshake(ctx kim.DialerContext) (net.Conn, error) {
	logger.Info("start tcp dial: ", ctx.Address)
	// 1 调用net.Dial拨号
	conn, err := net.DialTimeout("tcp", ctx.Address, ctx.Timeout)
	if err != nil {
		return nil, err
	}
	// 2. 发送用户认证信息，示例就是userid
	err = tcp.WriteFrame(conn, kim.OpBinary, []byte(ctx.Id))
	if err != nil {
		return nil, err
	}
	// 3. return conn
	return conn, nil
}

package mock

import (
	"errors"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/naming"
	"github.com/klintcheng/kim/tcp"
	"github.com/klintcheng/kim/websocket"
)

type ServerDemo struct{}

func (s *ServerDemo) Start(id, protocol, addr string) {
	go func() {
		_ = http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	var srv kim.Server
	service := &naming.DefaultService{
		Id:       id,
		Protocol: protocol,
	}
	if protocol == "ws" {
		srv = websocket.NewServer(addr, service)
	} else if protocol == "tcp" {
		srv = tcp.NewServer(addr, service)
	}

	handler := &ServerHandler{}

	srv.SetReadWait(time.Minute)
	srv.SetAcceptor(handler)
	srv.SetMessageListener(handler)
	srv.SetStateListener(handler)

	err := srv.Start()
	if err != nil {
		panic(err)
	}
}

// ServerHandler ServerHandler
type ServerHandler struct {
}

// Accept this connection
func (h *ServerHandler) Accept(conn kim.Conn, timeout time.Duration) (string, kim.Meta, error) {
	// 1. 读取：客户端发送的鉴权数据包
	frame, err := conn.ReadFrame()
	if err != nil {
		return "", nil, err
	}
	// 2. 解析：数据包内容就是userId
	userID := string(frame.GetPayload())
	// 3. 鉴权：这里只是为了示例做一个fake验证，非空
	if userID == "" {
		return "", nil, errors.New("user id is invalid")
	}
	logger.Infof("logined %s", userID)
	return userID, nil, nil
}

// Receive default listener
func (h *ServerHandler) Receive(ag kim.Agent, payload []byte) {
	logger.Infof("srv received %s", string(payload))
	_ = ag.Push([]byte("ok"))
}

// Disconnect default listener
func (h *ServerHandler) Disconnect(id string) error {
	logger.Infof("disconnect %s", id)
	return nil
}

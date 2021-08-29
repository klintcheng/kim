package websocket

import (
	"bufio"
	"net"

	"github.com/gobwas/ws"
	"github.com/klintcheng/kim"
)

// Server is a websocket implement of the Server
type Upgrader struct {
}

// NewServer NewServer
func NewServer(listen string, service kim.ServiceRegistration) kim.Server {
	return kim.NewServer(listen, service, new(Upgrader))
}

func (u *Upgrader) Name() string {
	return "websocket.Server"
}

func (u *Upgrader) Upgrade(rawconn net.Conn, rd *bufio.Reader, wr *bufio.Writer) (kim.Conn, error) {
	_, err := ws.Upgrade(rawconn)
	if err != nil {
		return nil, err
	}
	conn := NewConnWithRW(rawconn, rd, wr)
	return conn, nil
}

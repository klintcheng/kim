package websocket

import (
	"bufio"
	"net"

	"github.com/gobwas/ws"
	"github.com/klintcheng/kim"
)

type Frame struct {
	raw ws.Frame
}

func (f *Frame) SetOpCode(code kim.OpCode) {
	f.raw.Header.OpCode = ws.OpCode(code)
}

func (f *Frame) GetOpCode() kim.OpCode {
	return kim.OpCode(f.raw.Header.OpCode)
}

func (f *Frame) SetPayload(payload []byte) {
	f.raw.Payload = payload
}

func (f *Frame) GetPayload() []byte {
	if f.raw.Header.Masked {
		ws.Cipher(f.raw.Payload, f.raw.Header.Mask, 0)
	}
	f.raw.Header.Masked = false
	return f.raw.Payload
}

type WsConn struct {
	net.Conn
	rd *bufio.Reader
	wr *bufio.Writer
}

func NewConn(conn net.Conn) kim.Conn {
	return &WsConn{
		Conn: conn,
		rd:   bufio.NewReaderSize(conn, 4096),
		wr:   bufio.NewWriterSize(conn, 1024),
	}
}

func NewConnWithRW(conn net.Conn, rd *bufio.Reader, wr *bufio.Writer) *WsConn {
	return &WsConn{
		Conn: conn,
		rd:   rd,
		wr:   wr,
	}
}

func (c *WsConn) ReadFrame() (kim.Frame, error) {
	f, err := ws.ReadFrame(c.rd)
	if err != nil {
		return nil, err
	}
	return &Frame{raw: f}, nil
}

func (c *WsConn) WriteFrame(code kim.OpCode, payload []byte) error {
	f := ws.NewFrame(ws.OpCode(code), true, payload)
	return ws.WriteFrame(c.wr, f)
}

func (c *WsConn) Flush() error {
	return c.wr.Flush()
}

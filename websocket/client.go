package websocket

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
)

// ClientOptions ClientOptions
type ClientOptions struct {
	Heartbeat time.Duration //登录超时
	ReadWait  time.Duration //读超时
	WriteWait time.Duration //写超时
}

// Client is a websocket implement of the terminal
type Client struct {
	sync.Mutex
	kim.Dialer
	once    sync.Once
	id      string
	name    string
	conn    net.Conn
	state   int32
	options ClientOptions
	Meta    map[string]string
}

// NewClient NewClient
func NewClient(id, name string, opts ClientOptions) kim.Client {
	return NewClientWithProps(id, name, make(map[string]string), opts)
}

func NewClientWithProps(id, name string, meta map[string]string, opts ClientOptions) kim.Client {
	if opts.WriteWait == 0 {
		opts.WriteWait = kim.DefaultWriteWait
	}
	if opts.ReadWait == 0 {
		opts.ReadWait = kim.DefaultReadWait
	}

	cli := &Client{
		id:      id,
		name:    name,
		options: opts,
		Meta:    meta,
	}
	return cli
}

// Connect to server
func (c *Client) Connect(addr string) error {
	_, err := url.Parse(addr)
	if err != nil {
		return err
	}
	if !atomic.CompareAndSwapInt32(&c.state, 0, 1) {
		return fmt.Errorf("client has connected")
	}
	// step 1 拨号及握手
	conn, err := c.Dialer.DialAndHandshake(kim.DialerContext{
		Id:      c.id,
		Name:    c.name,
		Address: addr,
		Timeout: kim.DefaultLoginWait,
	})
	if err != nil {
		atomic.CompareAndSwapInt32(&c.state, 1, 0)
		return err
	}
	if conn == nil {
		return fmt.Errorf("conn is nil")
	}
	c.conn = conn

	if c.options.Heartbeat > 0 {
		go func() {
			err := c.heartbeatloop(conn)
			if err != nil {
				logger.Error("heartbeatloop stopped ", err)
			}
		}()
	}
	return nil
}

// SetDialer 设置握手逻辑
func (c *Client) SetDialer(dialer kim.Dialer) {
	c.Dialer = dialer
}

//Send data to connection
func (c *Client) Send(payload []byte) error {
	if atomic.LoadInt32(&c.state) == 0 {
		return fmt.Errorf("connection is nil")
	}
	c.Lock()
	defer c.Unlock()
	err := c.conn.SetWriteDeadline(time.Now().Add(c.options.WriteWait))
	if err != nil {
		return err
	}
	// 客户端消息需要使用MASK
	return wsutil.WriteClientMessage(c.conn, ws.OpBinary, payload)
}

// Close 关闭
func (c *Client) Close() {
	c.once.Do(func() {
		if c.conn == nil {
			return
		}
		// graceful close connection
		_ = wsutil.WriteClientMessage(c.conn, ws.OpClose, nil)

		c.conn.Close()
		atomic.CompareAndSwapInt32(&c.state, 1, 0)
	})
}

// Read a frame ,this function is not safely for concurrent
func (c *Client) Read() (kim.Frame, error) {
	if c.conn == nil {
		return nil, errors.New("connection is nil")
	}
	if c.options.Heartbeat > 0 {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.options.ReadWait))
	}
	frame, err := ws.ReadFrame(c.conn)
	if err != nil {
		return nil, err
	}
	if frame.Header.OpCode == ws.OpClose {
		return nil, errors.New("remote side close the channel")
	}
	return &Frame{
		raw: frame,
	}, nil
}

func (c *Client) heartbeatloop(conn net.Conn) error {
	tick := time.NewTicker(c.options.Heartbeat)
	for range tick.C {
		// 发送一个ping的心跳包给服务端
		if err := c.ping(conn); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ping(conn net.Conn) error {
	c.Lock()
	defer c.Unlock()
	err := conn.SetWriteDeadline(time.Now().Add(c.options.WriteWait))
	if err != nil {
		return err
	}
	logger.Tracef("%s send ping to server", c.id)
	return wsutil.WriteClientMessage(conn, ws.OpPing, nil)
}

// ID return id
func (c *Client) ServiceID() string {
	return c.id
}

// Name Name
func (c *Client) ServiceName() string {
	return c.name
}

func (c *Client) GetMeta() map[string]string {
	return c.Meta
}

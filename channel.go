package kim

import (
	"errors"
	"sync"
	"time"

	"github.com/klintcheng/kim/logger"
)

// Channel a websocket implement of channel
type ChannelImpl struct {
	sync.Mutex
	id string
	Conn
	writechan chan []byte
	once      sync.Once
	writeWait time.Duration
	readwait  time.Duration
	closed    *Event
}

// NewChannel 创建一个新的节点,如果 conf 为空，会使用一个默认的配置
func NewChannel(id string, conn Conn) Channel {
	log := logger.WithFields(logger.Fields{
		"module": "tcp_channel",
		"id":     id,
	})
	ch := &ChannelImpl{
		id:        id,
		Conn:      conn,
		writechan: make(chan []byte, 5),
		closed:    NewEvent(),
		writeWait: DefaultWriteWait, //default value
		readwait:  DefaultReadWait,
	}
	go func() {
		err := ch.writeloop()
		if err != nil {
			log.Warn(err)
		}
		ch.Conn.Close()
	}()
	return ch
}

func (ch *ChannelImpl) writeloop() error {
	for {
		select {
		case payload := <-ch.writechan:
			err := ch.WriteFrame(OpBinary, payload)
			if err != nil {
				return err
			}
			chanlen := len(ch.writechan)
			for i := 0; i < chanlen; i++ {
				payload = <-ch.writechan
				err := ch.WriteFrame(OpBinary, payload)
				if err != nil {
					return err
				}
			}
			err = ch.Conn.Flush()
			if err != nil {
				return err
			}
		case <-ch.closed.Done():
			return nil
		}
	}
}

// ID id
func (ch *ChannelImpl) ID() string { return ch.id }

// Send 异步写数据
func (ch *ChannelImpl) Push(payload []byte) error {
	if ch.closed.HasFired() {
		return errors.New("channel has closed")
	}
	// 异步写
	ch.writechan <- payload
	return nil
}

// overwrite Conn
func (ch *ChannelImpl) WriteFrame(code OpCode, payload []byte) error {
	ch.Lock()
	defer ch.Unlock()
	_ = ch.Conn.SetWriteDeadline(time.Now().Add(ch.writeWait))
	return ch.Conn.WriteFrame(code, payload)
}

// Close 关闭连接
func (ch *ChannelImpl) Close() error {
	ch.once.Do(func() {
		close(ch.writechan)
		ch.closed.Fire()
	})
	return nil
}

// SetWriteWait 设置写超时
func (ch *ChannelImpl) SetWriteWait(writeWait time.Duration) {
	if writeWait == 0 {
		return
	}
	ch.writeWait = writeWait
}
func (ch *ChannelImpl) SetReadWait(readwait time.Duration) {
	if readwait == 0 {
		return
	}
	ch.writeWait = readwait
}

func (ch *ChannelImpl) Readloop(lst MessageListener) error {
	for {
		_ = ch.SetReadDeadline(time.Now().Add(ch.readwait))

		frame, err := ch.ReadFrame()
		if err != nil {
			return err
		}
		if frame.GetOpCode() == OpClose {
			return errors.New("remote side close the channel")
		}
		if frame.GetOpCode() == OpPing {
			_ = ch.WriteFrame(OpPong, nil)
			continue
		}
		payload := frame.GetPayload()
		if len(payload) == 0 {
			continue
		}
		// TODO: Optimization point
		go lst.Receive(ch, payload)
	}
}

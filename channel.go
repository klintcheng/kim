package kim

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/klintcheng/kim/logger"
	"github.com/panjf2000/ants/v2"
)

// ChannelImpl is a websocket implement of channel
type ChannelImpl struct {
	sync.Mutex
	id string
	Conn
	writechan chan []byte
	once      sync.Once
	writeWait time.Duration
	readwait  time.Duration
	closed    *Event
	gpool     *ants.Pool
}

// NewChannel NewChannel
func NewChannel(id string, conn Conn, gpool *ants.Pool) Channel {
	ch := &ChannelImpl{
		id:        id,
		Conn:      conn,
		writechan: make(chan []byte, 10),
		closed:    NewEvent(),
		writeWait: DefaultWriteWait, //default value
		readwait:  DefaultReadWait,
		gpool:     gpool,
	}
	go func() {
		err := ch.writeloop()
		if err != nil {
			logger.WithFields(logger.Fields{
				"module": "ChannelImpl",
				"id":     id,
			}).Info(err)
		}
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
			_ = ch.Conn.SetWriteDeadline(time.Now().Add(ch.writeWait))
			err = ch.Flush()
			if err != nil {
				return err
			}
		case <-ch.closed.Done():
			return nil
		}
	}
}

// ID id simpling server
func (ch *ChannelImpl) ID() string { return ch.id }

// Send 异步写数据
func (ch *ChannelImpl) Push(payload []byte) error {
	if ch.closed.HasFired() {
		return fmt.Errorf("channel %s has closed", ch.id)
	}
	// 异步写
	ch.writechan <- payload
	return nil
}

// Close 关闭连接
func (ch *ChannelImpl) Close() error {
	ch.once.Do(func() {
		ch.closed.Fire()
		close(ch.writechan)
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
	ch.Lock()
	defer ch.Unlock()
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
			logger.WithFields(logger.Fields{
				"struct": "ChannelImpl",
				"func":   "Readloop",
				"id":     ch.id,
			}).Trace("recv a ping; resp with a pong")

			_ = ch.WriteFrame(OpPong, nil)
			continue
		}
		payload := frame.GetPayload()
		if len(payload) == 0 {
			continue
		}

		err = ch.gpool.Submit(func() {
			lst.Receive(ch, payload)
		})
		if err != nil {
			return err
		}
	}
}

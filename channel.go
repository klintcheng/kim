package kim

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/klintcheng/kim/logger"
	"github.com/panjf2000/ants/v2"
)

// ChannelImpl is a websocket implement of channel
type ChannelImpl struct {
	id string
	Conn
	writechan chan []byte
	writeWait time.Duration
	readwait  time.Duration
	gpool     *ants.Pool
	state     int32 // 0 init 1 start 2 closed
}

// NewChannel NewChannel
func NewChannel(id string, conn Conn, gpool *ants.Pool) Channel {
	ch := &ChannelImpl{
		id:        id,
		Conn:      conn,
		writechan: make(chan []byte, 5),
		writeWait: DefaultWriteWait, //default value
		readwait:  DefaultReadWait,
		gpool:     gpool,
		state:     0,
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
	for payload := range ch.writechan {
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
	}
	logger.Debugf("channel %s writeloop exited", ch.id)
	return nil
}

// ID id simpling server
func (ch *ChannelImpl) ID() string { return ch.id }

// Send 异步写数据
func (ch *ChannelImpl) Push(payload []byte) error {
	if atomic.LoadInt32(&ch.state) != 1 {
		return fmt.Errorf("channel %s has closed", ch.id)
	}
	// 异步写
	ch.writechan <- payload
	return nil
}

// Close 关闭连接
func (ch *ChannelImpl) Close() error {
	if !atomic.CompareAndSwapInt32(&ch.state, 1, 2) {
		return fmt.Errorf("channel has started")
	}
	close(ch.writechan)
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
	if !atomic.CompareAndSwapInt32(&ch.state, 0, 1) {
		return fmt.Errorf("channel has started")
	}
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

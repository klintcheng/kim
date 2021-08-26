package benchmark

import (
	"bytes"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/stretchr/testify/assert"
)

func Benchmark_Usertalk(b *testing.B) {
	cli1, err := dialer.Login(wsurl, "test1")
	assert.Nil(b, err)

	offline := true

	if offline {
		cli2, err := dialer.Login(wsurl, "test2")
		assert.Nil(b, err)

		go func() {
			for {
				_, err := cli2.Read()
				if err != nil {
					return
				}
			}
		}()
	}
	var lock sync.Mutex

	b.ReportAllocs()
	b.ResetTimer()
	t1 := time.Now()
	b.Logf("cpu %d", runtime.NumCPU())

	runtime.GOMAXPROCS(runtime.NumCPU())

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			p := pkt.New(wire.CommandChatUserTalk, pkt.WithDest("test2"))
			p.WriteBody(&pkt.MessageReq{
				Type: 1,
				Body: "hello world",
			})
			err := cli1.Send(pkt.Marshal(p))
			assert.Nil(b, err)

			lock.Lock()
			frame, _ := cli1.Read()
			lock.Unlock()

			if frame.GetOpCode() == kim.OpBinary {
				packet, err := pkt.MustReadLogicPkt(bytes.NewBuffer(frame.GetPayload()))
				assert.Nil(b, err)
				assert.Equal(b, pkt.Status_Success, packet.Header.Status)
			}
		}
	})
	b.Logf("cost %v", time.Since(t1))
}

func Benchmark_grouptalk(t *testing.B) {
	cli1, err := dialer.Login(wsurl, "test1")
	assert.Nil(t, err)

	// 创建群
	p := pkt.New(wire.CommandGroupCreate)
	p.WriteBody(&pkt.GroupCreateReq{
		Name:    "group1",
		Owner:   "test1",
		Members: []string{"test1", "test2", "test3", "test4", "test5"},
	})
	err = cli1.Send(pkt.Marshal(p))
	assert.Nil(t, err)
	// 读取返回信息
	ack, err := cli1.Read()
	assert.Nil(t, err)

	ackp, _ := pkt.MustReadLogicPkt(bytes.NewBuffer(ack.GetPayload()))
	assert.Equal(t, pkt.Status_Success, ackp.GetStatus())
	assert.Equal(t, wire.CommandGroupCreate, ackp.GetCommand())

	var createresp pkt.GroupCreateResp
	err = ackp.ReadBody(&createresp)
	assert.Nil(t, err)
	group := createresp.GetGroupId()
	assert.NotEmpty(t, group)
	if group == "" {
		return
	}
	// 登录
	cli2, err := dialer.Login(wsurl, "test2")
	assert.Nil(t, err)
	cli3, err := dialer.Login(wsurl, "test3")
	assert.Nil(t, err)
	t1 := time.Now()

	var lock sync.Mutex

	t.ReportAllocs()
	t.ResetTimer()
	t.Logf("cpu %d", runtime.NumCPU())

	go func() {
		for {
			_, _ = cli1.Read()
		}
	}()
	t.RunParallel(func(p *testing.PB) {
		for p.Next() {
			// 发送消息
			gtalk := pkt.New(wire.CommandChatGroupTalk, pkt.WithDest(group)).WriteBody(&pkt.MessageReq{
				Type: 1,
				Body: "hellogroup",
			})
			err = cli1.Send(pkt.Marshal(gtalk))
			assert.Nil(t, err)
			// 读取消息
			lock.Lock()
			notify1, _ := cli2.Read()
			notify2, _ := cli3.Read()
			lock.Unlock()

			n1, _ := pkt.MustReadLogicPkt(bytes.NewBuffer(notify1.GetPayload()))
			assert.Equal(t, wire.CommandChatGroupTalk, n1.GetCommand())
			var notify pkt.MessagePush
			_ = n1.ReadBody(&notify)
			assert.Equal(t, "hellogroup", notify.Body)
			assert.Equal(t, int32(1), notify.Type)
			assert.Empty(t, notify.Extra)
			assert.Greater(t, notify.SendTime, t1.UnixNano())
			assert.Greater(t, notify.MessageId, int64(10000))

			n2, _ := pkt.MustReadLogicPkt(bytes.NewBuffer(notify2.GetPayload()))
			_ = n2.ReadBody(&notify)
			assert.Equal(t, "hellogroup", notify.Body)
		}
	})

	t.Logf("cost %v", time.Since(t1))
}

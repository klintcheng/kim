package benchmark

import (
	"bytes"
	"fmt"
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

	if !offline {
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
	memberNums := 100
	var members = make([]string, memberNums)
	for i := 0; i < memberNums; i++ {
		members[i] = fmt.Sprintf("test%d", i+1)
	}
	// 创建群
	p := pkt.New(wire.CommandGroupCreate)

	p.WriteBody(&pkt.GroupCreateReq{
		Name:    "group1",
		Owner:   "test1",
		Members: members,
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

	onlines := memberNums / 2
	for i := onlines; i < memberNums; i++ {
		clix, err := dialer.Login(wsurl, fmt.Sprintf("test%d", i+1))
		assert.Nil(t, err)
		go func(cli kim.Client) {
			for {
				_, err := cli.Read()
				if err != nil {
					return
				}
			}
		}(clix)
	}
	t1 := time.Now()

	var lock sync.Mutex

	t.ReportAllocs()
	t.ResetTimer()
	t.Logf("cpu %d", runtime.NumCPU())

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
			_, _ = cli1.Read()
			lock.Unlock()
		}
	})

	t.Logf("cost %v", time.Since(t1))
}

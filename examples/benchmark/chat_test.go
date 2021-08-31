package benchmark

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/stretchr/testify/assert"
)

func Benchmark_Usertalk(b *testing.B) {
	cli1, err := dialer.Login(wsurl, "test1", appSecret)
	assert.Nil(b, err)

	offline := true

	if !offline {
		cli2, err := dialer.Login(wsurl, "test2", appSecret)
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

	b.ReportAllocs()
	b.ResetTimer()
	t1 := time.Now()

	for i := 0; i < b.N; i++ {
		p := pkt.New(wire.CommandChatUserTalk, pkt.WithDest("test2"))
		p.WriteBody(&pkt.MessageReq{
			Type: 1,
			Body: "hello world",
		})
		err := cli1.Send(pkt.Marshal(p))
		assert.Nil(b, err)
		frame, _ := cli1.Read()

		if frame.GetOpCode() == kim.OpBinary {
			packet, err := pkt.MustReadLogicPkt(bytes.NewBuffer(frame.GetPayload()))
			assert.Nil(b, err)
			assert.Equal(b, pkt.Status_Success, packet.Header.Status)
		}
	}
	b.Logf("cost %v", time.Since(t1))
}

func Benchmark_grouptalk(b *testing.B) {
	cli1, err := dialer.Login(wsurl, "test1", appSecret)
	assert.Nil(b, err)
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
	assert.Nil(b, err)
	// 读取返回信息
	ack, err := cli1.Read()
	assert.Nil(b, err)

	ackp, _ := pkt.MustReadLogicPkt(bytes.NewBuffer(ack.GetPayload()))
	assert.Equal(b, pkt.Status_Success, ackp.GetStatus())
	assert.Equal(b, wire.CommandGroupCreate, ackp.GetCommand())

	var createresp pkt.GroupCreateResp
	err = ackp.ReadBody(&createresp)
	assert.Nil(b, err)
	group := createresp.GetGroupId()
	assert.NotEmpty(b, group)
	if group == "" {
		return
	}

	onlines := memberNums / 2
	for i := 1; i < onlines; i++ {
		clix, err := dialer.Login(wsurl, fmt.Sprintf("test%d", i), appSecret)
		assert.Nil(b, err)
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

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 发送消息
		gtalk := pkt.New(wire.CommandChatGroupTalk, pkt.WithDest(group)).WriteBody(&pkt.MessageReq{
			Type: 1,
			Body: "hellogroup",
		})
		err = cli1.Send(pkt.Marshal(gtalk))
		assert.Nil(b, err)
		// 读取消息
		_, _ = cli1.Read()
	}

	b.Logf("cost %v", time.Since(t1))

	cli1.Close()
}

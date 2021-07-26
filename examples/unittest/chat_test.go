package unittest

import (
	"bytes"
	"testing"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/stretchr/testify/assert"
)

func Test_Usertalk(t *testing.T) {
	cli1, err := login("test1")
	assert.Nil(t, err)

	cli2, err := login("test2")
	assert.Nil(t, err)

	p := pkt.New(wire.CommandChatUserTalk, pkt.WithDest("test2"))
	p.WriteBody(&pkt.MessageReq{
		Type: 1,
		Body: "hello world",
	})
	err = cli1.Send(pkt.Marshal(p))
	assert.Nil(t, err)

	// resp
	frame, _ := cli1.Read()
	assert.Equal(t, kim.OpBinary, frame.GetOpCode())
	packet, err := pkt.MustReadLogicPkt(bytes.NewBuffer(frame.GetPayload()))
	assert.Nil(t, err)
	assert.Equal(t, pkt.Status_Success, packet.Header.Status)
	var resp pkt.MessageResp
	_ = packet.ReadBody(&resp)
	assert.Greater(t, resp.MessageId, int64(1000))
	assert.Greater(t, resp.SendTime, int64(1000))
	t.Log(&resp)

	// push message
	frame, err = cli2.Read()
	assert.Nil(t, err)
	packet, err = pkt.MustReadLogicPkt(bytes.NewBuffer(frame.GetPayload()))
	assert.Nil(t, err)
	var push pkt.MessagePush
	_ = packet.ReadBody(&push)
	assert.Equal(t, resp.MessageId, push.MessageId)
	assert.Equal(t, resp.SendTime, push.SendTime)
	assert.Equal(t, "hello world", push.Body)
	assert.Equal(t, int32(1), push.Type)
	t.Log(&push)
}

func Test_grouptalk(t *testing.T) {
	cli1, err := login("test1")
	assert.Nil(t, err)

	// 创建群
	p := pkt.New(wire.CommandGroupCreate)
	p.WriteBody(&pkt.GroupCreateReq{
		Name:    "group1",
		Owner:   "test1",
		Members: []string{"test1", "test2", "test3", "test4"},
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
	cli2, err := login("test2")
	assert.Nil(t, err)
	cli3, err := login("test3")
	assert.Nil(t, err)
	t1 := time.Now()

	// 发送消息
	gtalk := pkt.New(wire.CommandChatGroupTalk, pkt.WithDest(group)).WriteBody(&pkt.MessageReq{
		Type: 1,
		Body: "hellogroup",
	})
	err = cli1.Send(pkt.Marshal(gtalk))
	assert.Nil(t, err)

	ack, _ = cli1.Read()
	ackp, _ = pkt.MustReadLogicPkt(bytes.NewBuffer(ack.GetPayload()))
	assert.Equal(t, pkt.Status_Success, ackp.GetStatus())

	// 读取消息
	notify1, _ := cli2.Read()
	n1, _ := pkt.MustReadLogicPkt(bytes.NewBuffer(notify1.GetPayload()))
	assert.Equal(t, wire.CommandChatGroupTalk, n1.GetCommand())
	var notify pkt.MessagePush
	_ = n1.ReadBody(&notify)
	assert.Equal(t, "hellogroup", notify.Body)
	assert.Equal(t, int32(1), notify.Type)
	assert.Empty(t, notify.Extra)
	assert.Greater(t, notify.SendTime, t1.UnixNano())
	assert.Greater(t, notify.MessageId, int64(10000))

	notify2, _ := cli3.Read()
	n2, _ := pkt.MustReadLogicPkt(bytes.NewBuffer(notify2.GetPayload()))
	_ = n2.ReadBody(&notify)
	assert.Equal(t, "hellogroup", notify.Body)

	t.Logf("cost %v", time.Since(t1))
}

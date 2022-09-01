package unittest

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func Test_offline(t *testing.T) {
	src := fmt.Sprintf("u%d", time.Now().Unix())
	cli, err := dialer.Login(wsurl, src)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	dest := fmt.Sprintf("u%d", time.Now().Unix()+1)
	count := 10
	for i := 0; i < count; i++ {
		p := pkt.New(wire.CommandChatUserTalk, pkt.WithDest(dest))
		p.WriteBody(&pkt.MessageReq{
			Type: 1,
			Body: "hello world",
		})
		err := cli.Send(pkt.Marshal(p))
		if err != nil {
			logger.Error(err)
			return
		}
		// wait ack
		_, _ = cli.Read()
	}

	destcli, err := dialer.Login(wsurl, dest)
	assert.Nil(t, err)

	// request offline message index
	p := pkt.New(wire.CommandOfflineIndex)
	p.WriteBody(&pkt.MessageIndexReq{})
	_ = destcli.Send(pkt.Marshal(p))

	var indexResp pkt.MessageIndexResp
	err = Read(destcli, &indexResp)
	assert.Nil(t, err)

	assert.Equal(t, count, len(indexResp.Indexes))
	assert.Equal(t, src, indexResp.Indexes[0].AccountB)
	assert.Equal(t, int32(0), indexResp.Indexes[0].Direction)
	t.Log(indexResp.Indexes)

	var ids = make([]int64, count)
	for i, idx := range indexResp.Indexes {
		ids[i] = idx.MessageId
	}
	t.Log(ids)

	lastMessageId := ids[count-1]

	// read again
	p = pkt.New(wire.CommandOfflineIndex)
	p.WriteBody(&pkt.MessageIndexReq{
		MessageId: lastMessageId,
	})
	_ = destcli.Send(pkt.Marshal(p))

	var indexResp2 pkt.MessageIndexResp
	err = Read(destcli, &indexResp2)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(indexResp2.Indexes))

	// request offline message content
	p = pkt.New(wire.CommandOfflineContent)
	p.WriteBody(&pkt.MessageContentReq{
		MessageIds: ids,
	})
	_ = destcli.Send(pkt.Marshal(p))
	var contentResp pkt.MessageContentResp
	err = Read(destcli, &contentResp)
	assert.Nil(t, err)
	t.Log(contentResp.Contents)
	assert.Equal(t, count, len(contentResp.Contents))
	assert.Equal(t, "hello world", contentResp.Contents[0].Body)
	assert.Equal(t, int32(1), contentResp.Contents[0].Type)
}

func Read(cli kim.Client, body proto.Message) error {
	frame, err := cli.Read()
	if err != nil {
		return err
	}
	packet, _ := pkt.MustReadLogicPkt(bytes.NewBuffer(frame.GetPayload()))
	if packet.GetStatus() != pkt.Status_Success {
		return fmt.Errorf("received status :%v", packet.GetStatus())
	}
	return packet.ReadBody(body)
}

package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/klintcheng/kim/wire/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_Message(t *testing.T) {

	messageService := NewMessageService("http://localhost:8080")

	m := rpc.Message{
		Type: 1,
		Body: "hello world",
	}
	dest := fmt.Sprintf("u%d", time.Now().Unix())
	_, err := messageService.InsertUser(app, &rpc.InsertMessageReq{
		Sender:   "test1",
		Dest:     dest,
		SendTime: time.Now().UnixNano(),
		Message:  &m,
	})
	assert.Nil(t, err)

	resp, err := messageService.GetMessageIndex(app, &rpc.GetOfflineMessageIndexReq{
		Account: dest,
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp.List))

	index := resp.List[0]
	assert.Equal(t, "test1", index.AccountB)

	resp2, err := messageService.GetMessageContent(app, &rpc.GetOfflineMessageContentReq{
		MessageIds: []int64{index.MessageId},
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp2.List))
	content := resp2.List[0]
	assert.Equal(t, m.Body, content.Body)
	assert.Equal(t, m.Type, content.Type)
	assert.Equal(t, index.MessageId, content.Id)

	//again
	resp, err = messageService.GetMessageIndex(app, &rpc.GetOfflineMessageIndexReq{
		Account:   dest,
		MessageId: index.MessageId,
	})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(resp.List))

	resp, err = messageService.GetMessageIndex(app, &rpc.GetOfflineMessageIndexReq{
		Account: dest,
	})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(resp.List))
}

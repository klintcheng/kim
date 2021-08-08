package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/klintcheng/kim/wire/rpc"
	"github.com/stretchr/testify/assert"
)

var messageService = NewMessageServiceWithSRV("http", &resty.SRVRecord{
	Domain:  "consul",
	Service: "royal",
})

func Test_Message(t *testing.T) {

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

func Test_Group_Message(t *testing.T) {
	resp, err := groupService.Create(app, &rpc.CreateGroupReq{
		Name:    "test",
		Owner:   "test1",
		Members: []string{"test1", "test2", "test3"},
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, resp.GroupId)

	m := rpc.Message{
		Type: 1,
		Body: "hello world",
	}
	dest := resp.GroupId
	_, err = messageService.InsertGroup(app, &rpc.InsertMessageReq{
		Sender:   "test1",
		Dest:     dest,
		SendTime: time.Now().UnixNano(),
		Message:  &m,
	})
	assert.Nil(t, err)

	indexresp, err := messageService.GetMessageIndex(app, &rpc.GetOfflineMessageIndexReq{
		Account: "test1",
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(indexresp.List))
	assert.Equal(t, int32(1), indexresp.List[0].Direction)

	indexresp2, err := messageService.GetMessageIndex(app, &rpc.GetOfflineMessageIndexReq{
		Account: "test2",
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(indexresp2.List))
	assert.Equal(t, int32(0), indexresp2.List[0].Direction)
}

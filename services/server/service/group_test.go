package service

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/klintcheng/kim/wire/rpc"
	"github.com/stretchr/testify/assert"
)

const app = "kim_t"

var groupService = NewGroupServiceWithSRV("http", &resty.SRVRecord{
	Domain:  "consul",
	Service: "royal",
})

func TestGroupService(t *testing.T) {

	resp, err := groupService.Create(app, &rpc.CreateGroupReq{
		Name:    "test",
		Owner:   "test1",
		Members: []string{"test1", "test2"},
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, resp.GroupId)
	t.Log(resp.GroupId)

	mresp, err := groupService.Members(app, &rpc.GroupMembersReq{
		GroupId: resp.GroupId,
	})
	assert.Nil(t, err)

	assert.Equal(t, 2, len(mresp.Users))
	assert.Equal(t, "test1", mresp.Users[0].Account)
	assert.Equal(t, "test2", mresp.Users[1].Account)

	err = groupService.Join(app, &rpc.JoinGroupReq{
		Account: "test3",
		GroupId: resp.GroupId,
	})
	assert.Nil(t, err)

	mresp, err = groupService.Members(app, &rpc.GroupMembersReq{
		GroupId: resp.GroupId,
	})
	assert.Nil(t, err)

	assert.Equal(t, 3, len(mresp.Users))
	assert.Equal(t, "test3", mresp.Users[2].Account)
	assert.Equal(t, "test2", mresp.Users[1].Account)

	err = groupService.Quit(app, &rpc.QuitGroupReq{
		Account: "test2",
		GroupId: resp.GroupId,
	})
	assert.Nil(t, err)

	mresp, err = groupService.Members(app, &rpc.GroupMembersReq{
		GroupId: resp.GroupId,
	})
	assert.Nil(t, err)

	assert.Equal(t, 2, len(mresp.Users))
	assert.Equal(t, "test1", mresp.Users[0].Account)
	assert.Equal(t, "test3", mresp.Users[1].Account)
}

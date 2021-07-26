package service

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/wire/rpc"
	"google.golang.org/protobuf/proto"
)

type Group interface {
	Create(app string, req *rpc.CreateGroupReq) (*rpc.CreateGroupResp, error)
	Members(app string, req *rpc.GroupMembersReq) (*rpc.GroupMembersResp, error)
	Join(app string, req *rpc.JoinGroupReq) error
	Quit(app string, req *rpc.QuitGroupReq) error
}

type GroupHttp struct {
	url string
	cli *resty.Client
}

func NewGroupService(url string) Group {
	return &GroupHttp{
		url: url,
		cli: resty.New().SetRetryCount(3).SetTimeout(time.Second*5).SetHeader("userAgent", "kim_server"),
	}
}

func (m *GroupHttp) Create(app string, req *rpc.CreateGroupReq) (*rpc.CreateGroupResp, error) {
	path := fmt.Sprintf("%s/api/%s/group", m.url, app)

	body, _ := proto.Marshal(req)
	response, err := m.cli.R().SetHeader("Content-Type", "application/x-protobuf").SetBody(body).Post(path)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("GroupHttp.Create response.StatusCode() = %d, want 200", response.StatusCode())
	}
	var resp rpc.CreateGroupResp
	_ = proto.Unmarshal(response.Body(), &resp)
	logger.Debugf("GroupHttp.Create resp: %v", &resp)
	return &resp, nil
}

func (m *GroupHttp) Members(app string, req *rpc.GroupMembersReq) (*rpc.GroupMembersResp, error) {
	path := fmt.Sprintf("%s/api/%s/group/members/%s", m.url, app, req.GroupId)

	response, err := m.cli.R().SetHeader("Content-Type", "application/x-protobuf").Get(path)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("GroupHttp.Members response.StatusCode() = %d, want 200", response.StatusCode())
	}
	var resp rpc.GroupMembersResp
	_ = proto.Unmarshal(response.Body(), &resp)
	logger.Debugf("GroupHttp.Members resp: %v", &resp)
	return &resp, nil
}

func (m *GroupHttp) Join(app string, req *rpc.JoinGroupReq) error {
	path := fmt.Sprintf("%s/api/%s/group/member", m.url, app)
	body, _ := proto.Marshal(req)
	response, err := m.cli.R().SetHeader("Content-Type", "application/x-protobuf").SetBody(body).Post(path)
	if err != nil {
		return err
	}
	if response.StatusCode() != 200 {
		return fmt.Errorf("GroupHttp.Join response.StatusCode() = %d, want 200", response.StatusCode())
	}
	return nil
}

func (m *GroupHttp) Quit(app string, req *rpc.QuitGroupReq) error {
	path := fmt.Sprintf("%s/api/%s/group/member", m.url, app)
	body, _ := proto.Marshal(req)
	response, err := m.cli.R().SetHeader("Content-Type", "application/x-protobuf").SetBody(body).Delete(path)
	if err != nil {
		return err
	}
	if response.StatusCode() != 200 {
		return fmt.Errorf("GroupHttp.Quit response.StatusCode() = %d, want 200", response.StatusCode())
	}
	return nil
}

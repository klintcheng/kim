package handler

import (
	"errors"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/services/server/service"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/klintcheng/kim/wire/rpc"
)

var ErrNoDestination = errors.New("dest is empty")

type ChatHandler struct {
	msgService   service.Message
	groupService service.Group
}

func NewChatHandler(message service.Message, group service.Group) *ChatHandler {
	return &ChatHandler{
		msgService:   message,
		groupService: group,
	}
}

func (h *ChatHandler) DoUserTalk(ctx kim.Context) {
	// validate
	if ctx.Header().Dest == "" {
		_ = ctx.RespWithError(pkt.Status_NoDestination, ErrNoDestination)
		return
	}
	// 1. 解包
	var req pkt.MessageReq
	if err := ctx.ReadBody(&req); err != nil {
		_ = ctx.RespWithError(pkt.Status_InvalidPacketBody, err)
		return
	}
	// 2. 获取接收方的位置信息
	receiver := ctx.Header().GetDest()
	loc, err := ctx.GetLocation(receiver, "")
	if err != nil && err != kim.ErrSessionNil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}
	// 3. 保存离线消息
	sendTime := time.Now().UnixNano()
	resp, err := h.msgService.InsertUser(ctx.Session().GetApp(), &rpc.InsertMessageReq{
		Sender:   ctx.Session().GetAccount(),
		Dest:     receiver,
		SendTime: sendTime,
		Message: &rpc.Message{
			Type:  req.GetType(),
			Body:  req.GetBody(),
			Extra: req.GetExtra(),
		},
	})
	if err != nil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}
	msgId := resp.MessageId

	// 4. 如果接收方在线，就推送一条消息过去。
	if loc != nil {
		if err = ctx.Dispatch(&pkt.MessagePush{
			MessageId: msgId,
			Type:      req.GetType(),
			Body:      req.GetBody(),
			Extra:     req.GetExtra(),
			Sender:    ctx.Session().GetAccount(),
			SendTime:  sendTime,
		}, loc); err != nil {
			_ = ctx.RespWithError(pkt.Status_SystemException, err)
			return
		}
	}
	// 5. 返回一条resp消息
	_ = ctx.Resp(pkt.Status_Success, &pkt.MessageResp{
		MessageId: msgId,
		SendTime:  sendTime,
	})
}

func (h *ChatHandler) DoGroupTalk(ctx kim.Context) {
	if ctx.Header().GetDest() == "" {
		_ = ctx.RespWithError(pkt.Status_NoDestination, ErrNoDestination)
		return
	}
	var req pkt.MessageReq
	if err := ctx.ReadBody(&req); err != nil {
		_ = ctx.RespWithError(pkt.Status_InvalidPacketBody, err)
		return
	}
	group := ctx.Header().GetDest()
	sendTime := time.Now().UnixNano()

	membersResp, err := h.groupService.Members(ctx.Session().GetApp(), &rpc.GroupMembersReq{
		GroupId: group,
	})
	if err != nil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}
	var members = make([]string, len(membersResp.Users))
	for i, user := range membersResp.Users {
		members[i] = user.Account
	}
	// find group members location
	locs, err := ctx.GetLocations(members...)
	if err != nil && err != kim.ErrSessionNil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}
	resp, err := h.msgService.InsertGroup(ctx.Session().GetApp(), &rpc.InsertMessageReq{
		Sender:   ctx.Session().GetAccount(),
		Dest:     group,
		SendTime: sendTime,
		Message: &rpc.Message{
			Type:  req.GetType(),
			Body:  req.GetBody(),
			Extra: req.GetExtra(),
		},
	})
	if err != nil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}

	// push to receiver
	if len(locs) > 0 {
		if err = ctx.Dispatch(&pkt.MessagePush{
			MessageId: resp.MessageId,
			Type:      req.GetType(),
			Body:      req.GetBody(),
			Extra:     req.GetExtra(),
			Sender:    ctx.Session().GetAccount(),
			SendTime:  sendTime,
		}, locs...); err != nil {
			_ = ctx.RespWithError(pkt.Status_SystemException, err)
			return
		}
	}
	// resp
	_ = ctx.Resp(pkt.Status_Success, &pkt.MessageResp{
		MessageId: resp.MessageId,
		SendTime:  sendTime,
	})
}

func (h *ChatHandler) DoTalkAck(ctx kim.Context) {
	var req pkt.MessageAckReq
	if err := ctx.ReadBody(&req); err != nil {
		_ = ctx.RespWithError(pkt.Status_InvalidPacketBody, err)
		return
	}
	err := h.msgService.SetAck(ctx.Session().GetApp(), &rpc.AckMessageReq{
		Account:   ctx.Session().GetAccount(),
		MessageId: req.GetMessageId(),
	})
	if err != nil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}
	_ = ctx.Resp(pkt.Status_Success, nil)
}

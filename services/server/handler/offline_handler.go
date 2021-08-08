package handler

import (
	"errors"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/services/server/service"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/klintcheng/kim/wire/rpc"
)

type OfflineHandler struct {
	msgService service.Message
}

func NewOfflineHandler(message service.Message) *OfflineHandler {
	return &OfflineHandler{
		msgService: message,
	}
}

func (h *OfflineHandler) DoSyncIndex(ctx kim.Context) {
	var req pkt.MessageIndexReq
	if err := ctx.ReadBody(&req); err != nil {
		_ = ctx.RespWithError(pkt.Status_InvalidPacketBody, err)
		return
	}
	resp, err := h.msgService.GetMessageIndex(ctx.Session().GetApp(), &rpc.GetOfflineMessageIndexReq{
		Account:   ctx.Session().GetAccount(),
		MessageId: req.GetMessageId(),
	})
	if err != nil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}
	var list = make([]*pkt.MessageIndex, len(resp.List))
	for i, val := range resp.List {
		list[i] = &pkt.MessageIndex{
			MessageId: val.MessageId,
			Direction: val.Direction,
			SendTime:  val.SendTime,
			AccountB:  val.AccountB,
			Group:     val.Group,
		}
	}
	_ = ctx.Resp(pkt.Status_Success, &pkt.MessageIndexResp{
		Indexes: list,
	})
}

func (h *OfflineHandler) DoSyncContent(ctx kim.Context) {
	var req pkt.MessageContentReq
	if err := ctx.ReadBody(&req); err != nil {
		_ = ctx.RespWithError(pkt.Status_InvalidPacketBody, err)
		return
	}
	if len(req.MessageIds) == 0 {
		_ = ctx.RespWithError(pkt.Status_InvalidPacketBody, errors.New("empty MessageIds"))
		return
	}
	resp, err := h.msgService.GetMessageContent(ctx.Session().GetApp(), &rpc.GetOfflineMessageContentReq{
		MessageIds: req.MessageIds,
	})
	if err != nil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}
	var list = make([]*pkt.MessageContent, len(resp.List))
	for i, val := range resp.List {
		list[i] = &pkt.MessageContent{
			MessageId: val.Id,
			Type:      val.Type,
			Body:      val.Body,
			Extra:     val.Extra,
		}
	}
	_ = ctx.Resp(pkt.Status_Success, &pkt.MessageContentResp{
		Contents: list,
	})
}

package handler

import (
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/wire/pkt"
)

type LoginHandler struct {
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (h *LoginHandler) DoSysLogin(ctx kim.Context) {
	// 1. 序列化
	var session pkt.Session
	if err := ctx.ReadBody(&session); err != nil {
		_ = ctx.RespWithError(pkt.Status_InvalidPacketBody, err)
		return
	}

	logger.WithFields(logger.Fields{
		"Func":      "Login",
		"ChannelId": session.GetChannelId(),
		"Account":   session.GetAccount(),
		"RemoteIP":  session.GetRemoteIP(),
	}).Info("do login")

	// 2. 检查当前账号是否已经登陆在其它地方
	old, err := ctx.GetLocation(session.Account, "")
	if err != nil && err != kim.ErrSessionNil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}

	if old != nil {
		// 3. 通知这个用户下线
		_ = ctx.Dispatch(&pkt.KickoutNotify{
			ChannelId: old.ChannelId,
		}, old)
	}

	// 4. 添加到会话管理器中
	err = ctx.Add(&session)
	if err != nil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}
	// 5. 返回一个登陆成功的消息
	var resp = &pkt.LoginResp{
		ChannelId: session.ChannelId,
	}
	_ = ctx.Resp(pkt.Status_Success, resp)
}

func (h *LoginHandler) DoSysLogout(ctx kim.Context) {
	logger.WithFields(logger.Fields{
		"Func":      "Logout",
		"ChannelId": ctx.Session().GetChannelId(),
		"Account":   ctx.Session().GetAccount(),
	}).Info("do Logout ")

	err := ctx.Delete(ctx.Session().GetAccount(), ctx.Session().GetChannelId())
	if err != nil {
		_ = ctx.RespWithError(pkt.Status_SystemException, err)
		return
	}

	_ = ctx.Resp(pkt.Status_Success, nil)
}

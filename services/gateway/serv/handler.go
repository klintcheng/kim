package serv

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/container"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/klintcheng/kim/wire/token"
)

var log = logger.WithFields(logger.Fields{
	"service": "gateway",
	"pkg":     "serv",
})

// Handler Handler
type Handler struct {
	ServiceID string
	AppSecret string
}

// Accept this connection
func (h *Handler) Accept(conn kim.Conn, timeout time.Duration) (string, error) {
	// 1. 读取登陆包
	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	frame, err := conn.ReadFrame()
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(frame.GetPayload())
	req, err := pkt.MustReadLogicPkt(buf)
	if err != nil {
		log.Error(err)
		return "", err
	}
	// 2. 必须是登录包
	if req.Command != wire.CommandLoginSignIn {
		resp := pkt.NewFrom(&req.Header)
		resp.Status = pkt.Status_InvalidCommand
		_ = conn.WriteFrame(kim.OpBinary, pkt.Marshal(resp))
		return "", fmt.Errorf("must be a SignIn command")
	}

	// 3. 反序列化Body
	var login pkt.LoginReq
	err = req.ReadBody(&login)
	if err != nil {
		return "", err
	}
	secret := h.AppSecret
	if secret == "" {
		secret = token.DefaultSecret
	}
	// 4. 使用默认的DefaultSecret 解析token
	tk, err := token.Parse(secret, login.Token)
	if err != nil {
		// 5. 如果token无效，就返回SDK一个Unauthorized消息
		resp := pkt.NewFrom(&req.Header)
		resp.Status = pkt.Status_Unauthorized
		_ = conn.WriteFrame(kim.OpBinary, pkt.Marshal(resp))
		return "", err
	}
	// 6. 生成一个全局唯一的ChannelID
	id := generateChannelID(h.ServiceID, tk.Account)
	log.Infof("accept %v channel:%s", tk, id)

	req.ChannelId = id
	req.WriteBody(&pkt.Session{
		Account:   tk.Account,
		ChannelId: id,
		GateId:    h.ServiceID,
		App:       tk.App,
		RemoteIP:  getIP(conn.RemoteAddr().String()),
	})
	// 7. 把login.转发给Login服务
	err = container.Forward(wire.SNLogin, req)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Receive default listener
func (h *Handler) Receive(ag kim.Agent, payload []byte) {
	buf := bytes.NewBuffer(payload)
	packet, err := pkt.Read(buf)
	if err != nil {
		log.Error(err)
		return
	}
	if basicPkt, ok := packet.(*pkt.BasicPkt); ok {
		if basicPkt.Code == pkt.CodePing {
			_ = ag.Push(pkt.Marshal(&pkt.BasicPkt{Code: pkt.CodePong}))
		}
		return
	}
	if logicPkt, ok := packet.(*pkt.LogicPkt); ok {
		logicPkt.ChannelId = ag.ID()

		err = container.Forward(logicPkt.ServiceName(), logicPkt)
		if err != nil {
			logger.WithFields(logger.Fields{
				"module": "handler",
				"id":     ag.ID(),
				"cmd":    logicPkt.Command,
				"dest":   logicPkt.Dest,
			}).Error(err)
		}
	}

}

// Disconnect default listener
func (h *Handler) Disconnect(id string) error {
	log.Infof("disconnect %s", id)

	logout := pkt.New(wire.CommandLoginSignOut, pkt.WithChannel(id))
	err := container.Forward(wire.SNLogin, logout)
	if err != nil {
		logger.WithFields(logger.Fields{
			"module": "handler",
			"id":     id,
		}).Error(err)
	}
	return nil
}

var ipExp = regexp.MustCompile(string("\\:[0-9]+$"))

func getIP(remoteAddr string) string {
	if remoteAddr == "" {
		return ""
	}
	return ipExp.ReplaceAllString(remoteAddr, "")
}

func generateChannelID(serviceID, account string) string {
	return fmt.Sprintf("%s_%s_%d", serviceID, account, wire.Seq.Next())
}

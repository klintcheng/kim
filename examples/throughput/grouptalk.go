package throughput

import (
	"bytes"
	"fmt"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
)

func grouptalk(wsurl, appSecret string, count int, memberCount int, onlinePercent float32) error {
	cli1, err := dialer.Login(wsurl, "test1", appSecret)
	if err != nil {
		return err
	}
	var members = make([]string, memberCount)
	for i := 0; i < memberCount; i++ {
		members[i] = fmt.Sprintf("test%d", i+1)
	}
	// 创建群
	p := pkt.New(wire.CommandGroupCreate)
	p.WriteBody(&pkt.GroupCreateReq{
		Name:    "group1",
		Owner:   "test1",
		Members: members,
	})
	if err = cli1.Send(pkt.Marshal(p)); err != nil {
		return err
	}
	// 读取返回信息
	ack, _ := cli1.Read()
	ackp, _ := pkt.MustReadLogicPkt(bytes.NewBuffer(ack.GetPayload()))
	if pkt.Status_Success != ackp.GetStatus() {
		return fmt.Errorf("create group failed")
	}

	var createresp pkt.GroupCreateResp
	_ = ackp.ReadBody(&createresp)
	group := createresp.GetGroupId()

	onlines := int(float32(memberCount) * onlinePercent)
	if onlines < 1 {
		onlines = 1
	}
	for i := 1; i < onlines; i++ {
		clix, err := dialer.Login(wsurl, fmt.Sprintf("test%d", i), appSecret)
		if err != nil {
			return err
		}
		go func(cli kim.Client) {
			for {
				_, err := cli.Read()
				if err != nil {
					return
				}
			}
		}(clix)
	}
	t1 := time.Now()
	go func() {
		for i := 0; i < count; i++ {
			gtalk := pkt.New(wire.CommandChatGroupTalk, pkt.WithDest(group)).WriteBody(&pkt.MessageReq{
				Type: 1,
				Body: "hellogroup",
			})
			err := cli1.Send(pkt.Marshal(gtalk))
			if err != nil {
				logger.Error(err)
				return
			}
		}
	}()

	for i := 0; i < count; i++ {
		_, err := cli1.Read()
		if err != nil {
			return err
		}
	}

	dur := time.Since(t1)
	logger.Infof("send message count %d; cost time: %v; qps:%v", count, dur, int64(count*1000)/dur.Milliseconds())
	return nil
}

package throughput

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/report"
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

	r := report.New(os.Stdout, count)
	t1 := time.Now()
	requests := make(map[uint32]time.Time, count)
	go func() {
		for i := 0; i < count; i++ {
			gtalk := pkt.New(wire.CommandChatGroupTalk, pkt.WithDest(group)).WriteBody(&pkt.MessageReq{
				Type: 1,
				Body: "hellogroup",
			})
			requests[gtalk.Sequence] = time.Now()

			err := cli1.Send(pkt.Marshal(gtalk))
			if err != nil {
				logger.Error(err)
				return
			}
		}
	}()

	for i := 0; i < count; i++ {
		frame, err := cli1.Read()
		if err != nil {
			r.Add(&report.Result{
				Err:           err,
				ContentLength: 11,
			})
			continue
		}
		ack, err := pkt.MustReadLogicPkt(bytes.NewBuffer(frame.GetPayload()))
		r.Add(&report.Result{
			Duration:      time.Since(requests[ack.GetSequence()]),
			Err:           err,
			ContentLength: 11,
			StatusCode:    int(ack.GetStatus()),
		})
	}

	r.Finalize(time.Since(t1))
	return nil
}

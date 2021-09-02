package kimbench

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/report"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/panjf2000/ants/v2"
)

func grouptalk(wsurl, appSecret string, threads, count int, memberCount int, onlinePercent float32) error {
	cli1, err := dialer.Login(wsurl, "test1", appSecret)
	if err != nil {
		return err
	}
	var members = make([]string, memberCount)
	for i := 0; i < memberCount; i++ {
		members[i] = fmt.Sprintf("test_%d", i+1)
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
		clix, err := dialer.Login(wsurl, fmt.Sprintf("test_%d", i), appSecret)
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

	clis, err := loginMulti(wsurl, appSecret, 2, threads)
	if err != nil {
		return err
	}

	pool, _ := ants.NewPool(threads, ants.WithPreAlloc(true))
	defer pool.Release()

	r := report.New(os.Stdout, count)
	t1 := time.Now()

	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		cli := clis[i%threads]
		_ = pool.Submit(func() {
			defer func() {
				wg.Done()
			}()

			t0 := time.Now()
			p := pkt.New(wire.CommandChatGroupTalk, pkt.WithDest(group))
			p.WriteBody(&pkt.MessageReq{
				Type: 1,
				Body: "hello world",
			})
			// 发送消息
			err := cli.Send(pkt.Marshal(p))
			if err != nil {
				r.Add(&report.Result{
					Err:           err,
					ContentLength: 11,
				})
				return
			}
			// 读取Resp消息
			_, err = cli.Read()
			if err != nil {
				r.Add(&report.Result{
					Err:           err,
					ContentLength: 11,
				})
				return
			}
			r.Add(&report.Result{
				Duration:   time.Since(t0),
				Err:        err,
				StatusCode: 0,
			})
		})
	}

	wg.Wait()
	r.Finalize(time.Since(t1))
	return nil
}

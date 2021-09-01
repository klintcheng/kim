package throughput

import (
	"bytes"
	"os"
	"time"

	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/report"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
)

func usertalk(wsurl, appSecret string, count int, online bool) error {
	cli1, err := dialer.Login(wsurl, "test1", appSecret)
	if err != nil {
		return err
	}

	if online {
		cli2, _ := dialer.Login(wsurl, "test2")

		go func() {
			for {
				_, err := cli2.Read()
				if err != nil {
					return
				}
			}
		}()
	}
	r := report.New(os.Stdout, count)
	t1 := time.Now()

	requests := make(map[uint32]time.Time, count)
	go func() {
		for i := 0; i < count; i++ {
			p := pkt.New(wire.CommandChatUserTalk, pkt.WithDest("test2"))
			requests[p.Sequence] = time.Now()

			p.WriteBody(&pkt.MessageReq{
				Type: 1,
				Body: "hello world",
			})
			err := cli1.Send(pkt.Marshal(p))
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

package throughput

import (
	"time"

	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
)

func usertalk(wsurl, appSecret string, count int, offline bool) error {
	cli1, err := dialer.Login(wsurl, "test1", appSecret)
	if err != nil {
		return err
	}

	if !offline {
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
	t1 := time.Now()
	go func() {
		for i := 0; i < count; i++ {
			p := pkt.New(wire.CommandChatUserTalk, pkt.WithDest("test2"))
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
		_, err := cli1.Read()
		if err != nil {
			return err
		}
	}
	dur := time.Since(t1)
	logger.Infof("cost count %d cost time: %v qps:%v", count, dur, int64(count*1000)/dur.Milliseconds())
	return nil
}

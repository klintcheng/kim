package throughput

import (
	"fmt"
	"sync"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/logger"
	"github.com/panjf2000/ants/v2"
)

func login(wsurl, appSecret string, count int, keep time.Duration) error {
	p, _ := ants.NewPool(count, ants.WithPreAlloc(true))
	defer p.Release()

	var wg sync.WaitGroup
	wg.Add(count)
	t1 := time.Now()
	clis := make([]kim.Client, count)
	for i := 0; i < count; i++ {
		idx := i
		_ = p.Submit(func() {
			cli, err := dialer.Login(wsurl, fmt.Sprintf("test%d", idx+1), appSecret)
			if err != nil {
				logger.Error(err)
			} else {
				clis[idx] = cli
			}
			wg.Done()
		})
	}
	wg.Wait()
	dur := time.Since(t1)
	logger.Infof("cost count %d cost time: %v qps:%v", count, dur, int64(count*1000)/dur.Milliseconds())

	logger.Infof("keep login for %v", keep)
	for _, cli := range clis {
		cli.Close()
	}
	logger.Infoln("shutdown..")
	return nil
}

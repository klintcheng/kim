package throughput

import (
	"fmt"
	"sync"
	"time"

	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/logger"
	"github.com/panjf2000/ants/v2"
)

func login(wsurl, appSecret string, count int) error {
	t1 := time.Now()
	p, _ := ants.NewPool(count, ants.WithPreAlloc(true))
	defer p.Release()
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		_ = p.Submit(func() {
			_, err := dialer.Login(wsurl, fmt.Sprintf("test%d", i+1), appSecret)
			if err != nil {
				logger.Error(err)
			}
			wg.Done()
		})
	}
	wg.Wait()
	dur := time.Since(t1)
	logger.Infof("cost count %d cost time: %v qps:%v", count, dur, int64(count*1000)/dur.Milliseconds())
	return nil
}

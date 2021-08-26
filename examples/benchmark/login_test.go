package benchmark

import (
	"fmt"
	"testing"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/websocket"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

const wsurl = "ws://localhost:8000"

func login(account string) (kim.Client, error) {
	cli := websocket.NewClient(account, "unittest", websocket.ClientOptions{})

	cli.SetDialer(&dialer.ClientDialer{})
	err := cli.Connect(wsurl)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func Benchmark_Login(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	t0 := time.Now()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := login(fmt.Sprintf("user_%v", ksuid.New()))
			assert.Nil(b, err)
		}
	})

	b.Logf("logined %d cost %v", b.N, time.Since(t0).Milliseconds())
}

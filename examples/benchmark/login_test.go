package benchmark

import (
	"fmt"
	"testing"
	"time"

	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/websocket"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func Benchmark_Login(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	t0 := time.Now()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			cli := websocket.NewClient(fmt.Sprintf("user_%v", ksuid.New()), "echo", websocket.ClientOptions{
				Heartbeat: time.Second * 30,
				ReadWait:  time.Minute * 3,
				WriteWait: time.Second * 10,
			})

			cli.SetDialer(&dialer.ClientDialer{})
			err := cli.Connect("ws://localhost:8000")
			assert.Nil(b, err)
		}
	})

	b.Logf("logined %d cost %v", b.N, time.Since(t0).Milliseconds())
}

package benchmark

import (
	"fmt"
	"testing"
	"time"

	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/wire/token"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

// const wsurl = "ws://119.3.4.216:8000"
// const appSecret = "kingimcloud.cn"

const wsurl = "ws://localhost:8000"
const appSecret = token.DefaultSecret

func Benchmark_Login(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	t0 := time.Now()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := dialer.Login(wsurl, fmt.Sprintf("user_%v", ksuid.New()), appSecret)
			assert.Nil(b, err)
		}
	})

	b.Logf("logined %d cost %v", b.N, time.Since(t0))
}

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

const wsurl = "ws://124.71.204.19:8000"
const appSecret = token.DefaultSecret

func Benchmark_Login(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	t0 := time.Now()
	for i := 0; i < b.N; i++ {
		_, err := dialer.Login(wsurl, fmt.Sprintf("user_%v", ksuid.New()), appSecret)
		assert.Nil(b, err)
	}

	b.Logf("logined %d cost %v", b.N, time.Since(t0))
}

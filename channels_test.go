package kim

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/segmentio/ksuid"
)

func Benchmark_ChannelsAdd(b *testing.B) {
	ctrl := gomock.NewController(b)

	chs := NewChannels(10)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			ch := NewMockChannel(ctrl)
			id := ksuid.New().String()
			ch.EXPECT().ID().AnyTimes().Return(id)
			chs.Add(ch)
			chs.Get(id)
		}
	})
}

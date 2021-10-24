package serv

import (
	"testing"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/naming"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestRouteSelector_Lookup(t *testing.T) {

	srvs := []kim.Service{
		&naming.DefaultService{
			Id:   "s1",
			Meta: map[string]string{"zone": "zone_ali_01"},
		},
		&naming.DefaultService{
			Id:   "s2",
			Meta: map[string]string{"zone": "zone_ali_01"},
		},
		&naming.DefaultService{
			Id:   "s3",
			Meta: map[string]string{"zone": "zone_ali_01"},
		},
		&naming.DefaultService{
			Id:   "s4",
			Meta: map[string]string{"zone": "zone_ali_02"},
		},
		&naming.DefaultService{
			Id:   "s5",
			Meta: map[string]string{"zone": "zone_ali_02"},
		},
		&naming.DefaultService{
			Id:   "s6",
			Meta: map[string]string{"zone": "zone_ali_03"},
		},
	}

	rs, err := NewRouteSelector("../route.json")
	assert.Nil(t, err)

	packet := pkt.New(wire.CommandChatUserTalk, pkt.WithChannel(ksuid.New().String()))
	packet.AddStringMeta(MetaKeyApp, "kim")
	packet.AddStringMeta(MetaKeyAccount, "test1")
	hit := rs.Lookup(&packet.Header, srvs)
	assert.Equal(t, "s6", hit)

	hits := make(map[string]int)
	for i := 0; i < 100; i++ {
		header := pkt.Header{
			ChannelId: ksuid.New().String(),
			Meta: []*pkt.Meta{
				{
					Type:  pkt.MetaType_string,
					Key:   MetaKeyApp,
					Value: ksuid.New().String(),
				},
				{
					Type:  pkt.MetaType_string,
					Key:   MetaKeyAccount,
					Value: ksuid.New().String(),
				},
			},
		}
		hit = rs.Lookup(&header, srvs)
		hits[hit]++
	}
	t.Log(hits)
}

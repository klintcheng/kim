package apis

import (
	"testing"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/naming"
	"github.com/klintcheng/kim/services/router/conf"
	"github.com/stretchr/testify/assert"
)

func Test_selectIdc(t *testing.T) {
	got := selectIdc("test1", &conf.Region{
		Idcs: []conf.IDC{
			{ID: "SH_ALI"},
			{ID: "HZ_ALI"},
			{ID: "SH_TENCENT"},
		},
		Slots: []byte{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2},
	})
	assert.NotNil(t, got)
	t.Log(got)
}

func Test_selectGateways(t *testing.T) {
	got := selectGateways("test11", []kim.ServiceRegistration{
		&naming.DefaultService{Id: "g1"},
		&naming.DefaultService{Id: "g2"},
	}, 3)
	assert.Equal(t, len(got), 2)

	got = selectGateways("test11", []kim.ServiceRegistration{
		&naming.DefaultService{Id: "g1"},
		&naming.DefaultService{Id: "g2"},
		&naming.DefaultService{Id: "g3"},
	}, 3)
	assert.Equal(t, len(got), 3)

	got = selectGateways("test11", []kim.ServiceRegistration{
		&naming.DefaultService{Id: "g1"},
		&naming.DefaultService{Id: "g2"},
		&naming.DefaultService{Id: "g3"},
		&naming.DefaultService{Id: "g4"},
		&naming.DefaultService{Id: "g5"},
		&naming.DefaultService{Id: "g6"},
	}, 3)

	t.Log(got)
	assert.Equal(t, len(got), 3)
}

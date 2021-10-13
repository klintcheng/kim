package router

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/klintcheng/kim/services/router/apis"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func Test_Lookup(t *testing.T) {
	cli := resty.New()
	cli.SetHeader("Content-Type", "application/json")
	shDomain := "ws://kingimcloud.com"
	hzDomain := "ws://kingimcloud2.com"

	shHit, hzHit := int(0), int(0)

	for i := 0; i < 1000; i++ {
		url := fmt.Sprintf("http://localhost:8100/api/lookup/%s", ksuid.New().String())

		var res apis.LookUpResp
		resp, err := cli.R().SetResult(&res).Get(url)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		if len(res.Domains) == 1 {
			if res.Domains[0] == shDomain {
				shHit++
			} else if res.Domains[0] == hzDomain {
				hzHit++
			}
		}
	}

	t.Logf("shHit %d ;hzHit %d", shHit, hzHit)
}

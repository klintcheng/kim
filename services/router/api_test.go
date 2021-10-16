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

	domains := make(map[string]int)
	for i := 0; i < 1000; i++ {
		url := fmt.Sprintf("http://localhost:8100/api/lookup/%s", ksuid.New().String())

		var res apis.LookUpResp
		resp, err := cli.R().SetResult(&res).Get(url)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		if len(res.Domains) > 0 {
			domain := res.Domains[0]
			domains[domain]++
		}
	}
	for domain, hit := range domains {
		fmt.Printf("domain: %s ;hit count: %d\n", domain, hit)
	}
}

package unittest

import (
	"testing"
	"time"

	"github.com/klintcheng/kim/examples/dialer"
	"github.com/stretchr/testify/assert"
)

// const wsurl = "ws://119.3.4.216:8000"
const wsurl = "ws://localhost:8000"

func Test_login(t *testing.T) {
	cli, err := dialer.Login(wsurl, "test1")
	assert.Nil(t, err)
	time.Sleep(time.Second * 2)
	cli.Close()
}

package unittest

import (
	"testing"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/examples/dialer"
	"github.com/klintcheng/kim/websocket"
	"github.com/stretchr/testify/assert"
)

func login(account string) (kim.Client, error) {
	cli := websocket.NewClient(account, "unittest", websocket.ClientOptions{})

	cli.SetDialer(&dialer.ClientDialer{})
	err := cli.Connect("ws://localhost:8000")
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func Test_login(t *testing.T) {
	cli, err := login("test1")
	assert.Nil(t, err)
	time.Sleep(time.Second * 2)
	cli.Close()
}

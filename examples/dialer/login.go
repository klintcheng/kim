package dialer

import (
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/websocket"
)

func Login(wsurl, account string, appSecrets ...string) (kim.Client, error) {
	cli := websocket.NewClient(account, "unittest", websocket.ClientOptions{
		Heartbeat: kim.DefaultHeartbeat,
	})
	secret := ""
	if len(appSecrets) > 0 {
		secret = appSecrets[0]
	}
	cli.SetDialer(&ClientDialer{
		AppSecret: secret,
	})
	err := cli.Connect(wsurl)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

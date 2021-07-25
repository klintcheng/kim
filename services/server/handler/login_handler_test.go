package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler_DoSysLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dispather := kim.NewMockDispather(ctrl)
	cache := kim.NewMockSessionStorage(ctrl)
	session := &pkt.Session{
		ChannelId: "channel1",
		Account:   "test1",
		GateId:    "gateway1",
	}
	// resp
	dispather.EXPECT().Push(session.GateId, []string{"channel1"}, gomock.Any()).Times(1)
	// kickout notify
	dispather.EXPECT().Push(session.GateId, []string{"channel2"}, gomock.Any()).Times(1)

	cache.EXPECT().GetLocation(session.Account, "").DoAndReturn(func(account string, device string) (*kim.Location, error) {
		return &kim.Location{
			ChannelId: "channel2",
			GateId:    "gateway1",
		}, kim.ErrSessionNil
	})

	cache.EXPECT().Add(gomock.Any()).Times(1).DoAndReturn(func(add *pkt.Session) error {
		assert.Equal(t, session.ChannelId, add.ChannelId)
		assert.Equal(t, session.Account, add.Account)
		return nil
	})

	loginreq := pkt.New(wire.CommandLoginSignIn).WriteBody(session)

	r := kim.NewRouter()
	// login
	loginHandler := NewLoginHandler()
	r.Handle(wire.CommandLoginSignIn, loginHandler.DoSysLogin)

	err := r.Serve(loginreq, dispather, cache, session)
	assert.Nil(t, err)
}

package kim

import (
	"bytes"
	"errors"

	"github.com/klintcheng/kim/wire/endian"
)

type Location struct {
	ChannelId string
	GateId    string
}

func (loc *Location) Bytes() []byte {
	if loc == nil {
		return []byte{}
	}
	buf := new(bytes.Buffer)
	_ = endian.WriteShortBytes(buf, []byte(loc.ChannelId))
	_ = endian.WriteShortBytes(buf, []byte(loc.GateId))
	return buf.Bytes()
}

func (loc *Location) Unmarshal(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("data is empty")
	}
	buf := bytes.NewBuffer(data)
	loc.ChannelId, err = endian.ReadShortString(buf)
	if err != nil {
		return
	}
	loc.GateId, err = endian.ReadShortString(buf)
	if err != nil {
		return
	}
	return
}

package pkt

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/klintcheng/kim/wire"
)

type Packet interface {
	Decode(r io.Reader) error
	Encode(w io.Writer) error
}

func MustReadLogicPkt(r io.Reader) (*LogicPkt, error) {
	val, err := Read(r)
	if err != nil {
		return nil, err
	}
	if lp, ok := val.(*LogicPkt); ok {
		return lp, nil
	}
	return nil, fmt.Errorf("packet is not a logic packet")
}

func MustReadBasicPkt(r io.Reader) (*BasicPkt, error) {
	val, err := Read(r)
	if err != nil {
		return nil, err
	}
	if bp, ok := val.(*BasicPkt); ok {
		return bp, nil
	}
	return nil, fmt.Errorf("packet is not a basic packet")
}

func Read(r io.Reader) (interface{}, error) {
	magic := wire.Magic{}
	_, err := io.ReadFull(r, magic[:])
	if err != nil {
		return nil, err
	}
	switch magic {
	case wire.MagicLogicPkt:
		p := new(LogicPkt)
		if err := p.Decode(r); err != nil {
			return nil, err
		}
		return p, nil
	case wire.MagicBasicPkt:
		p := new(BasicPkt)
		if err := p.Decode(r); err != nil {
			return nil, err
		}
		return p, nil
	default:
		return nil, fmt.Errorf("magic code %s is incorrect", magic)
	}
}

func Marshal(p Packet) []byte {
	buf := new(bytes.Buffer)
	kind := reflect.TypeOf(p).Elem()

	if kind.AssignableTo(reflect.TypeOf(LogicPkt{})) {
		_, _ = buf.Write(wire.MagicLogicPkt[:])
	} else if kind.AssignableTo(reflect.TypeOf(BasicPkt{})) {
		_, _ = buf.Write(wire.MagicBasicPkt[:])
	}
	_ = p.Encode(buf)
	return buf.Bytes()
}

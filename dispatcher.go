package kim

import "github.com/klintcheng/kim/wire/pkt"

// Dispatcher defined a component how a message be dispatched to gateway
type Dispatcher interface {
	Push(gateway string, channels []string, p *pkt.LogicPkt) error
}

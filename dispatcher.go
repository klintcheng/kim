package kim

import "github.com/klintcheng/kim/wire/pkt"

// Dispather defined a component how a message be dispatched to gateway
type Dispather interface {
	Push(gateway string, channels []string, p *pkt.LogicPkt) error
}

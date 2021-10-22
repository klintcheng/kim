package serv

import (
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/wire/pkt"
)

// MetaSelector MetaSelector
type MetaSelector struct {
}

// Lookup a server
func (s *MetaSelector) Lookup(header *pkt.Header, srvs []kim.Service) string {
	// ll := len(srvs)

	return srvs[0].ServiceID()
}

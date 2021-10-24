package serv

import (
	"hash/crc32"
	"math/rand"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/services/gateway/conf"
	"github.com/klintcheng/kim/wire/pkt"
)

// RouteSelector RouteSelector
type RouteSelector struct {
	route *conf.Route
}

func NewRouteSelector(configPath string) (*RouteSelector, error) {
	route, err := conf.ReadRoute(configPath)
	if err != nil {
		return nil, err
	}
	return &RouteSelector{
		route: route,
	}, nil
}

// Lookup a server
func (s *RouteSelector) Lookup(header *pkt.Header, srvs []kim.Service) string {
	app, _ := pkt.FindMeta(header.Meta, MetaKeyApp)
	account, _ := pkt.FindMeta(header.Meta, MetaKeyAccount)
	if app == nil || account == nil {
		ri := rand.Intn(len(srvs))
		return srvs[ri].ServiceID()
	}
	log := logger.WithFields(logger.Fields{
		"app":     app,
		"account": account,
	})

	// 1. 判断是否命中白名单
	zone, ok := s.route.Whitelist[app.(string)]
	if !ok {
		var key string
		switch s.route.RouteBy {
		case MetaKeyApp:
			key = app.(string)
		case MetaKeyAccount:
			key = account.(string)
		default:
			key = account.(string)
		}
		// 2. 未命中，通过权重计算出zone
		slot := hashcode(key) % len(s.route.Slots)
		i := s.route.Slots[slot]
		zone = s.route.Zones[i].ID
	} else {
		log.Infoln("hit a zone in whitelist", zone)
	}
	// 3. 过滤出当前zone的servers
	zoneSrvs := filterSrvs(srvs, zone)
	if len(zoneSrvs) == 0 {
		log.Warnf("select a random service from all due to no service found in zone %s", zone)
		ri := rand.Intn(len(srvs))
		return srvs[ri].ServiceID()
	}
	// 4. 选中一个服务
	srv := selectSrvs(zoneSrvs, account.(string))
	return srv.ServiceID()
}

func filterSrvs(srvs []kim.Service, zone string) []kim.Service {
	var res = make([]kim.Service, 0, len(srvs))
	for _, srv := range srvs {
		if zone == srv.GetMeta()["zone"] {
			res = append(res, srv)
		}
	}
	return res
}

func selectSrvs(srvs []kim.Service, account string) kim.Service {
	slots := make([]int, 0, len(srvs)*10)
	for i := range srvs {
		for j := 0; j < 10; j++ {
			slots = append(slots, i)
		}
	}
	slot := hashcode(account) % len(slots)
	return srvs[slots[slot]]
}

func hashcode(key string) int {
	hash32 := crc32.NewIEEE()
	hash32.Write([]byte(key))
	return int(hash32.Sum32())
}

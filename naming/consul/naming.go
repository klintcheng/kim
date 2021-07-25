package consul

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/naming"
)

const (
	KeyProtocol  = "protocol"
	KeyHealthURL = "health_url"
)

type Watch struct {
	Service   string
	Callback  func([]kim.ServiceRegistration)
	WaitIndex uint64
	Quit      chan struct{}
}

type Naming struct {
	sync.RWMutex
	cli    *api.Client
	watchs map[string]*Watch
}

func NewNaming(consulUrl string) (naming.Naming, error) {
	conf := api.DefaultConfig()
	conf.Address = consulUrl
	cli, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	naming := &Naming{
		cli:    cli,
		watchs: make(map[string]*Watch, 1),
	}

	return naming, nil
}

func (n *Naming) Find(name string, tags ...string) ([]kim.ServiceRegistration, error) {
	services, _, err := n.load(name, 0, tags...)
	if err != nil {
		return nil, err
	}
	return services, nil
}

// refresh service registration
func (n *Naming) load(name string, waitIndex uint64, tags ...string) ([]kim.ServiceRegistration, *api.QueryMeta, error) {
	opts := &api.QueryOptions{
		UseCache:  true,
		MaxAge:    time.Minute, // MaxAge limits how old a cached value will be returned if UseCache is true.
		WaitIndex: waitIndex,
	}
	catalogServices, meta, err := n.cli.Catalog().ServiceMultipleTags(name, tags, opts)
	if err != nil {
		return nil, meta, err
	}

	services := make([]kim.ServiceRegistration, len(catalogServices))
	for i, s := range catalogServices {
		if s.Checks.AggregatedStatus() != api.HealthPassing {
			logger.Debugf("load service: id:%s name:%s %s:%d Status:%s", s.ServiceID, s.ServiceName, s.ServiceAddress, s.ServicePort, s.Checks.AggregatedStatus())
			continue
		}
		services[i] = &naming.DefaultService{
			Id:       s.ServiceID,
			Name:     s.ServiceName,
			Address:  s.ServiceAddress,
			Port:     s.ServicePort,
			Protocol: s.ServiceMeta[KeyProtocol],
			Tags:     s.ServiceTags,
			Meta:     s.ServiceMeta,
		}
	}
	logger.Debugf("load service: %v, meta:%v", services, meta)
	return services, meta, nil
}

func (n *Naming) Register(s kim.ServiceRegistration) error {
	reg := &api.AgentServiceRegistration{
		ID:      s.ServiceID(),
		Name:    s.ServiceName(),
		Address: s.PublicAddress(),
		Port:    s.PublicPort(),
		Tags:    s.GetTags(),
		Meta:    s.GetMeta(),
	}
	if reg.Meta == nil {
		reg.Meta = make(map[string]string)
	}
	reg.Meta[KeyProtocol] = s.GetProtocol()

	// consul健康检查
	healthURL := s.GetMeta()[KeyHealthURL]
	if healthURL != "" {
		check := new(api.AgentServiceCheck)
		check.CheckID = fmt.Sprintf("%s_normal", s.ServiceID())
		check.HTTP = healthURL
		check.Timeout = "1s" // http timeout
		check.Interval = "10s"
		check.DeregisterCriticalServiceAfter = "20s"
		reg.Check = check
	}

	err := n.cli.Agent().ServiceRegister(reg)
	return err
}

func (n *Naming) Deregister(serviceID string) error {
	return n.cli.Agent().ServiceDeregister(serviceID)
}

func (n *Naming) Subscribe(serviceName string, callback func([]kim.ServiceRegistration)) error {
	n.Lock()
	defer n.Unlock()
	if _, ok := n.watchs[serviceName]; ok {
		return errors.New("serviceName has already been registered")
	}
	w := &Watch{
		Service:  serviceName,
		Callback: callback,
		Quit:     make(chan struct{}, 1),
	}
	n.watchs[serviceName] = w

	go n.watch(w)
	return nil
}

func (n *Naming) Unsubscribe(serviceName string) error {
	n.Lock()
	defer n.Unlock()
	wh, ok := n.watchs[serviceName]

	delete(n.watchs, serviceName)
	if ok {
		close(wh.Quit)
	}
	return nil
}

func (n *Naming) watch(wh *Watch) {
	stopped := false

	var doWatch = func(service string, callback func([]kim.ServiceRegistration)) {
		services, meta, err := n.load(service, wh.WaitIndex) // <-- blocking untill services has changed
		if err != nil {
			logger.Warn(err)
			return
		}
		select {
		case <-wh.Quit:
			stopped = true
			logger.Infof("watch %s stopped", wh.Service)
			return
		default:
		}

		wh.WaitIndex = meta.LastIndex
		if callback != nil {
			callback(services)
		}
	}

	// build WaitIndex
	doWatch(wh.Service, nil)
	for !stopped {
		doWatch(wh.Service, wh.Callback)
	}
}

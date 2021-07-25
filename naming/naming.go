package naming

import (
	"errors"

	"github.com/klintcheng/kim"
)

// errors
var (
	ErrNotFound = errors.New("service no found")
)

// Naming defined methods of the naming service
type Naming interface {
	Find(serviceName string, tags ...string) ([]kim.ServiceRegistration, error)
	Subscribe(serviceName string, callback func(services []kim.ServiceRegistration)) error
	Unsubscribe(serviceName string) error
	Register(service kim.ServiceRegistration) error
	Deregister(serviceID string) error
}

package naming

import (
	"fmt"
)

//Service define a Service
type ServiceRegistration interface {
	ServiceID() string
	ServiceName() string
	//ip or doamin
	PublicAddress() string
	PublicPort() int
	DialURL() string
	GetProtocol() string
	GetNamespace() string
	GetTags() []string
	GetMeta() map[string]string
	// SetTags(tags []string)
	// SetMeta(meta map[string]string)
	String() string
}

// DefaultService Service Impl
type DefaultService struct {
	Id        string
	Name      string
	Address   string
	Port      int
	Protocol  string
	Namespace string
	Tags      []string
	Meta      map[string]string
}

// NewEntry NewEntry
func NewEntry(id, name, protocol string, address string, port int) ServiceRegistration {
	return &DefaultService{
		Id:       id,
		Name:     name,
		Address:  address,
		Port:     port,
		Protocol: protocol,
	}
}

// ID returns the ServiceImpl ID
func (e *DefaultService) ServiceID() string {
	return e.Id
}

// Namespace Namespace
func (e *DefaultService) GetNamespace() string { return e.Namespace }

// Name Name
func (e *DefaultService) ServiceName() string { return e.Name }

// Address Address
func (e *DefaultService) PublicAddress() string { return e.Address }

func (e *DefaultService) PublicPort() int { return e.Port }

// Protocol Protocol
func (e *DefaultService) GetProtocol() string { return e.Protocol }

func (e *DefaultService) DialURL() string {
	if e.Protocol == "tcp" {
		return fmt.Sprintf("%s:%d", e.Address, e.Port)
	}
	return fmt.Sprintf("%s://%s:%d", e.Protocol, e.Address, e.Port)
}

// Tags Tags
func (e *DefaultService) GetTags() []string { return e.Tags }

// Meta Meta
func (e *DefaultService) GetMeta() map[string]string { return e.Meta }

// // SetTags SetTags
// func (e *DefaultService) SetTags(tags []string) {
// 	e.tags = tags
// }

// // SetMeta SetMeta
// func (e *DefaultService) SetMeta(meta map[string]string) {
// 	e.meta = meta
// }

func (e *DefaultService) String() string {
	return fmt.Sprintf("Id:%s,Name:%s,Address:%s,Port:%d,Ns:%s,Tags:%v,Meta:%v", e.Id, e.Name, e.Address, e.Port, e.Namespace, e.Tags, e.Meta)
}

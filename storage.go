package kim

import (
	"errors"

	"github.com/klintcheng/kim/wire/pkt"
)

// ErrNil
var ErrSessionNil = errors.New("err:session nil")

// SessionStorage defined a session storage which provides based functions as save,delete,find a session
type SessionStorage interface {
	// Add a session
	Add(session *pkt.Session) error
	// Delete a session
	Delete(account string, channelId string) error
	// Get session by channelId
	Get(channelId string) (*pkt.Session, error)
	// Get Locations by accounts
	GetLocations(account ...string) ([]*Location, error)
	// Get Location by account and device
	GetLocation(account string, device string) (*Location, error)
}

package kim

import (
	"sync"

	"github.com/klintcheng/kim/logger"
)

// ChannelMap ChannelMap
type ChannelMap interface {
	Add(channel Channel)
	Remove(id string)
	Get(id string) (channel Channel, ok bool)
	All() []Channel
}

// ChannelsImpl ChannelMap
type ChannelsImpl struct {
	channels *sync.Map
}

// NewChannels NewChannels
func NewChannels(num int) ChannelMap {
	return &ChannelsImpl{
		channels: new(sync.Map),
	}
}

// Add addChannel
func (ch *ChannelsImpl) Add(channel Channel) {
	if channel.ID() == "" {
		logger.WithFields(logger.Fields{
			"module": "ChannelsImpl",
		}).Error("channel id is required")
	}

	ch.channels.Store(channel.ID(), channel)
}

// Remove addChannel
func (ch *ChannelsImpl) Remove(id string) {
	ch.channels.Delete(id)
}

// Get Get
func (ch *ChannelsImpl) Get(id string) (Channel, bool) {
	if id == "" {
		logger.WithFields(logger.Fields{
			"module": "ChannelsImpl",
		}).Error("channel id is required")
	}

	val, ok := ch.channels.Load(id)
	if !ok {
		return nil, false
	}
	return val.(Channel), true
}

// All return channels
func (ch *ChannelsImpl) All() []Channel {
	arr := make([]Channel, 0)
	ch.channels.Range(func(key, val interface{}) bool {
		arr = append(arr, val.(Channel))
		return true
	})
	return arr
}

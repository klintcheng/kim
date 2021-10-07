package storage

import (
	"time"

	"github.com/chasex/redis-go-cluster"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/wire/pkt"
	"google.golang.org/protobuf/proto"
)

type RedisClusterStorage struct {
	cli redis.Cluster
}

func NewRedisClusterStorage(cli redis.Cluster) kim.SessionStorage {
	return &RedisClusterStorage{
		cli: cli,
	}
}

func (r *RedisClusterStorage) Add(session *pkt.Session) error {
	// save kim.Location
	loc := kim.Location{
		ChannelId: session.ChannelId,
		GateId:    session.GateId,
	}
	ex := int(LocationExpired / time.Second)
	locKey := KeyLocation(session.Account, "")
	_, err := r.cli.Do("SET", locKey, loc.Bytes(), "EX", ex)
	if err != nil {
		return err
	}
	// save session
	snKey := KeySession(session.ChannelId)
	buf, _ := proto.Marshal(session)
	_, err = r.cli.Do("SET", snKey, buf, "EX", ex)
	if err != nil {
		return err
	}
	return nil
}

// Delete a session
func (r *RedisClusterStorage) Delete(account string, channelId string) error {
	locKey := KeyLocation(account, "")
	_, err := r.cli.Do("DEL", locKey)
	if err != nil {
		return err
	}

	snKey := KeySession(channelId)
	_, err = r.cli.Do("DEL", snKey)
	if err != nil {
		return err
	}
	return nil
}

// GetByID get session by sessionID
func (r *RedisClusterStorage) Get(ChannelId string) (*pkt.Session, error) {
	snKey := KeySession(ChannelId)
	bts, err := redis.Bytes(r.cli.Do("GET", snKey))
	if err != nil {
		if err == redis.ErrNil {
			return nil, kim.ErrSessionNil
		}
		return nil, err
	}
	var session pkt.Session
	_ = proto.Unmarshal(bts, &session)
	return &session, nil
}

func (r *RedisClusterStorage) GetLocations(accounts ...string) ([]*kim.Location, error) {
	keys := KeyLocations(accounts...)
	list, err := redis.Values(r.cli.Do("MGET", keys))
	if err != nil {
		return nil, err
	}
	var result = make([]*kim.Location, 0)
	for _, l := range list {
		if l == nil {
			continue
		}
		var loc kim.Location
		_ = loc.Unmarshal([]byte(l.(string)))
		result = append(result, &loc)
	}
	if len(result) == 0 {
		return nil, kim.ErrSessionNil
	}
	return result, nil
}

func (r *RedisClusterStorage) GetLocation(account string, device string) (*kim.Location, error) {
	key := KeyLocation(account, device)
	bts, err := redis.Bytes(r.cli.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, kim.ErrSessionNil
		}
		return nil, err
	}
	var loc kim.Location
	_ = loc.Unmarshal(bts)
	return &loc, nil
}

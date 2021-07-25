package storage

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/wire/pkt"
	"google.golang.org/protobuf/proto"
)

const (
	LocationExpired = time.Hour * 48
)

type RedisStorage struct {
	cli *redis.Client
}

func NewRedisStorage(cli *redis.Client) kim.SessionStorage {
	return &RedisStorage{
		cli: cli,
	}
}

func (r *RedisStorage) Add(session *pkt.Session) error {
	// save kim.Location
	loc := kim.Location{
		ChannelId: session.ChannelId,
		GateId:    session.GateId,
	}
	locKey := KeyLocation(session.Account, "")
	err := r.cli.Set(locKey, loc.Bytes(), LocationExpired).Err()
	if err != nil {
		return err
	}
	// save session
	snKey := KeySession(session.ChannelId)
	buf, _ := proto.Marshal(session)
	err = r.cli.Set(snKey, buf, LocationExpired).Err()
	if err != nil {
		return err
	}
	return nil
}

// Delete a session
func (r *RedisStorage) Delete(account string, channelId string) error {
	locKey := KeyLocation(account, "")
	err := r.cli.Del(locKey).Err()
	if err != nil {
		return err
	}

	snKey := KeySession(channelId)
	err = r.cli.Del(snKey).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetByID get session by sessionID
func (r *RedisStorage) Get(ChannelId string) (*pkt.Session, error) {
	snKey := KeySession(ChannelId)
	bts, err := r.cli.Get(snKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, kim.ErrSessionNil
		}
		return nil, err
	}
	var session pkt.Session
	_ = proto.Unmarshal(bts, &session)
	return &session, nil
}

func (r *RedisStorage) GetLocations(accounts ...string) ([]*kim.Location, error) {
	keys := KeyLocations(accounts...)
	list, err := r.cli.MGet(keys...).Result()
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

func (r *RedisStorage) GetLocation(account string, device string) (*kim.Location, error) {
	key := KeyLocation(account, device)
	bts, err := r.cli.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, kim.ErrSessionNil
		}
		return nil, err
	}
	var loc kim.Location
	_ = loc.Unmarshal(bts)
	return &loc, nil
}

func KeySession(channel string) string {
	return fmt.Sprintf("login:sn:%s", channel)
}

func KeyLocation(account, device string) string {
	if device == "" {
		return fmt.Sprintf("login:loc:%s", account)
	}
	return fmt.Sprintf("login:loc:%s:%s", account, device)
}

func KeyLocations(accounts ...string) []string {
	arr := make([]string, len(accounts))
	for i, account := range accounts {
		arr[i] = KeyLocation(account, "")
	}
	return arr
}

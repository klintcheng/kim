package storage

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/wire/pkt"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func Test_crud(t *testing.T) {
	cli, err := InitRedis("localhost:6379", "")
	assert.Nil(t, err)
	cc := NewRedisStorage(cli)
	err = cc.Add(&pkt.Session{
		ChannelId: "ch1",
		GateId:    "gateway1",
		Account:   "test1",
		Device:    "Phone",
	})
	assert.Nil(t, err)

	_ = cc.Add(&pkt.Session{
		ChannelId: "ch2",
		GateId:    "gateway1",
		Account:   "test2",
		Device:    "Pc",
	})

	session, err := cc.Get("ch1")
	assert.Nil(t, err)
	t.Log(session)
	assert.Equal(t, "ch1", session.ChannelId)
	assert.Equal(t, "gateway1", session.GateId)
	assert.Equal(t, "test1", session.Account)

	arr, err := cc.GetLocations("test1", "test2")
	assert.Nil(t, err)
	t.Log(arr)
	loc := arr[1]

	arr, err = cc.GetLocations("test6")
	assert.Equal(t, kim.ErrSessionNil, err)
	assert.Equal(t, 0, len(arr))

	assert.Equal(t, "ch2", loc.ChannelId)
	assert.Equal(t, "gateway1", loc.GateId)
}

func Benchmark_MGET(b *testing.B) {
	cli, err := InitRedis("localhost:6379", "")
	assert.Nil(b, err)
	cc := NewRedisStorage(cli)
	count := 100
	accounts := make([]string, count)
	for i := 0; i < 100; i++ {
		accounts[i] = fmt.Sprintf("account_%d", i)
		err = cc.Add(&pkt.Session{
			ChannelId: ksuid.New().String(),
			GateId:    "gateway1",
			Account:   accounts[i],
		})
		assert.Nil(b, err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := cc.GetLocations(accounts...)
			assert.Nil(b, err)
		}
	})
}

func Benchmark_getLocation(b *testing.B) {
	cli, err := InitRedis("localhost:6379", "")
	assert.Nil(b, err)
	cc := NewRedisStorage(cli)

	accs := make([]string, 100)
	for i := 0; i < 100; i++ {
		accs[i] = ksuid.New().String()
		_ = cc.Add(&pkt.Session{
			ChannelId: ksuid.New().String(),
			GateId:    "127_0_0_1_gateway1",
			Account:   accs[i],
			Zone:      "testtesttesttest",
			Isp:       "moblie",
			RemoteIP:  "127.0.0.1",
			App:       "kim",
			Tags:      []string{"tag1", "tag2"},
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			account := accs[rand.Intn(100)]
			_, err := cc.GetLocation(account, "")
			assert.Nil(b, err)
		}
	})
}

func Benchmark_getSession(b *testing.B) {
	cli, err := InitRedis("localhost:6379", "")
	assert.Nil(b, err)
	cc := NewRedisStorage(cli)

	ids := make([]string, 100)
	for i := 0; i < 100; i++ {
		ids[i] = ksuid.New().String()
		_ = cc.Add(&pkt.Session{
			ChannelId: ids[i],
			GateId:    "127_0_0_1_gateway1",
			Account:   ksuid.New().String(),
			Zone:      "testtesttesttest",
			Isp:       "moblie",
			RemoteIP:  "127.0.0.1",
			App:       "kim",
			Tags:      []string{"tag1", "tag2"},
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := cc.Get(ids[rand.Intn(100)])
			assert.Nil(b, err)
		}
	})
}

func InitRedis(addr string, pass string) (*redis.Client, error) {
	redisdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pass,
		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	})

	_, err := redisdb.Ping().Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return redisdb, nil
}

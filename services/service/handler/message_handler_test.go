package handler

import (
	"fmt"
	"testing"
	"time"

	"github.com/klintcheng/kim/services/service/database"
	"github.com/klintcheng/kim/wire/rpc"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var handler ServiceHandler

func init() {
	baseDb, _ := database.InitDb("mysql", "root:123456@tcp(127.0.0.1:3306)/kim_base?charset=utf8mb4&parseTime=True&loc=Local")
	messageDb, _ := database.InitDb("mysql", "root:123456@tcp(127.0.0.1:3306)/kim_message?charset=utf8mb4&parseTime=True&loc=Local")
	idgen, _ := database.NewIDGenerator(1)
	handler = ServiceHandler{
		MessageDb: messageDb,
		BaseDb:    baseDb,
		Idgen:     idgen,
	}
}

func Benchmark_InsertUserMessage(b *testing.B) {

	b.ResetTimer()
	b.SetBytes(1024)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = handler.insertUserMessage(&rpc.InsertMessageReq{
				Sender:   "test1",
				Dest:     ksuid.New().String(),
				SendTime: time.Now().UnixNano(),
				Message: &rpc.Message{
					Type: 1,
					Body: "hello",
				},
			})
		}
	})
}

func Benchmark_InsertGroup10Message(b *testing.B) {
	memberCount := 10

	var members = make([]string, memberCount)
	for i := 0; i < memberCount; i++ {
		members[i] = fmt.Sprintf("test%d", i+1)
	}

	groupId, err := handler.groupCreate(&rpc.CreateGroupReq{
		App:     "kim_t",
		Name:    "testg",
		Owner:   "test1",
		Members: members,
	})
	assert.Nil(b, err)

	b.ResetTimer()
	b.SetBytes(1024)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = handler.insertGroupMessage(&rpc.InsertMessageReq{
				Sender:   "test1",
				Dest:     groupId.Base36(),
				SendTime: time.Now().UnixNano(),
				Message: &rpc.Message{
					Type: 1,
					Body: "hello",
				},
			})
		}
	})
}

func Benchmark_InsertGroup50Message(b *testing.B) {
	memberCount := 50

	var members = make([]string, memberCount)
	for i := 0; i < memberCount; i++ {
		members[i] = fmt.Sprintf("test%d", i+1)
	}

	groupId, err := handler.groupCreate(&rpc.CreateGroupReq{
		App:     "kim_t",
		Name:    "testg",
		Owner:   "test1",
		Members: members,
	})
	assert.Nil(b, err)

	b.ResetTimer()
	b.SetBytes(1024)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = handler.insertGroupMessage(&rpc.InsertMessageReq{
				Sender:   "test1",
				Dest:     groupId.Base36(),
				SendTime: time.Now().UnixNano(),
				Message: &rpc.Message{
					Type: 1,
					Body: "hello",
				},
			})
		}
	})
}

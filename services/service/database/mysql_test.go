package database

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB
var idgen *IDGenerator

func init() {
	db, _ = InitMysqlDb("root:123456@tcp(127.0.0.1:3306)/kim_message?charset=utf8mb4&parseTime=True&loc=Local")

	_ = db.AutoMigrate(&MessageIndex{})
	_ = db.AutoMigrate(&MessageContent{})

	idgen, _ = NewIDGenerator(1)
}

func Benchmark_insert(b *testing.B) {

	sendTime := time.Now().UnixNano()
	b.ResetTimer()
	b.SetBytes(1024)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			idxs := make([]MessageIndex, 100)
			cid := idgen.Next().Int64()
			for i := 0; i < len(idxs); i++ {
				idxs[i] = MessageIndex{
					ID:        idgen.Next().Int64(),
					AccountA:  fmt.Sprintf("test_%d", cid),
					AccountB:  fmt.Sprintf("test_%d", i),
					SendTime:  sendTime,
					MessageID: cid,
				}
			}
			db.Create(&idxs)
		}
	})
}

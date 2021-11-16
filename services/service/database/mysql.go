package database

import (

	// just init

	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitMysqlDb init mysql database
func InitDb(driver string, dsn string) (*gorm.DB, error) {
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	defaultLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	var dialector gorm.Dialector
	if driver == "mysql" {
		dialector = mysql.Open(dsn)
	}
	// else if driver == "sqlite" {
	// dialector = sqlite.Open(dsn)
	// }

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: defaultLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",                              // table name prefix, table for `User` would be `t_users`
			SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		}})

	return db, err
}

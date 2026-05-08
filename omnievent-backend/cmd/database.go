package cmd

import (
	"fmt"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var DB *xorm.Engine

func InitDatabase() error {
	cfg := GetConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	var err error
	DB, err = xorm.NewEngine("postgres", dsn)
	if err != nil {
		return err
	}

	DB.SetLogLevel(log.LOG_INFO)
	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(100)

	return nil
}

func GetDB() *xorm.Engine {
	return DB
}

func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

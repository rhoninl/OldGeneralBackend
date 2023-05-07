package database

import (
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	rdb *redis.Client
)

func setupSql() *gorm.DB {
	var (
		user         = os.Getenv("DB_USER")
		password     = os.Getenv("DB_PASSWORD")
		address      = os.Getenv("DB_ADDRESS")
		port         = os.Getenv("DB_PORT")
		databaseName = os.Getenv("DB_NAME")
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, address, port, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Error to connect to database, error: %s", err.Error())
	}
	return db
}

func setupRedis() *redis.Client {
	opt := &redis.Options{
		Username: os.Getenv("Redis_USER"),
		Password: os.Getenv("Redis_PASSWORD"),
		Addr:     os.Getenv("Redis_ADDRESS"),
		DB:       0,
	}
	return redis.NewClient(opt)
}

func GetDB() *gorm.DB {
	if db == nil {
		db = setupSql()
	}
	return db
}

func GetRDB() *redis.Client {
	if rdb == nil {
		rdb = setupRedis()
	}
	return rdb
}

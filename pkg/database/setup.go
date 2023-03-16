package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

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

func GetDB() *gorm.DB {
	if db == nil {
		db = setupSql()
	}

	return db
}

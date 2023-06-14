package config

import (
	"fmt"
	"log"
	"time"

	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// Read and Write Connection
	dsnRW := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "user", "pass", "host", "3306", "database_name")
	db, err := gorm.Open(mysql.Open(dsnRW), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Minute)

	if err := db.Raw(helper.ConvertToInLineQuery("SET GLOBAL FOREIGN_KEY_CHECKS = 0;")).Error; err != nil {
		log.Fatal(err.Error())
	}

	return db
}

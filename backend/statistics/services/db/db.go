package db

import (
	"time"
	"ylem_statistics/config"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var db *gorm.DB

func Instance() (*gorm.DB, error) {
	if db == nil {
		var err error
		db, err = newInstance()
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func newInstance() (*gorm.DB, error) {
	log.Debug("Creating a new DB connection instance with DSN " + config.Cfg().DB.DSN)

	db, err := gorm.Open(clickhouse.Open(config.Cfg().DB.DSN), &gorm.Config{})
	if err != nil {
		log.Debug("DB connection creation failed: " + err.Error())
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Debug("DB connection creation failed: " + err.Error())
		return nil, err
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(300)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Debug("New DB connection created")

	return db, nil
}

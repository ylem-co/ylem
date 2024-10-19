package helpers

import (
	"fmt"
	"testing"
	"ylem_api/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormInstance *gorm.DB

func init() {
	if testing.Testing() {
		return
	}

	cfg := config.Cfg()

	dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true&loc=UTC",
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.Name);

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			logrus.New(),
			logger.Config{
				IgnoreRecordNotFoundError: true,
			}),
	})
	if err != nil {
		panic(err)
	}
	gormInstance = db
}

func GormInstance() *gorm.DB {
	return gormInstance
}

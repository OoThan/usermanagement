package ds

import (
	"log"

	"github.com/OoThan/usermanagement/config"
	"github.com/OoThan/usermanagement/internal/model"
	"github.com/OoThan/usermanagement/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func LoadDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.MysqlDNS()), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		// Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	logger.Sugar.Info("Successfully connected to MySQL")

	// migrate DB
	err = db.AutoMigrate(
		&model.User{},
		&model.UserLog{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitTestDB() {
	var err error
	db, err := gorm.Open(mysql.Open("sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the in-memory database:", err)
	}

	db.AutoMigrate(
		&model.User{},
		&model.UserLog{},
	)
}

package ds

import (
	"github.com/OoThan/usermanagement/config"
	"github.com/OoThan/usermanagement/internal/model"
	"github.com/OoThan/usermanagement/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

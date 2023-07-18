package model

import (
	"time"

	"gorm.io/gorm"
)

type UserLog struct {
	Id     uint64 `gorm:"column:id;primaryKey" json:"id"`
	UserId uint64 `gorm:"column:user_id"`

	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type UserLogMongo struct {
	UserID     string `bson:"user_id"`
	Event      string `bson:"event"`
	Data       User   `bson:"data"`
	Timestamps int64  `bson:"timestamps"`
}

package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id       uint64 `gorm:"column:id;primaryKey" json:"id"`
	Name     string `gorm:"column:name;unique;not null" json:"name"`
	Email    string `gorm:"column:email;unique;not null" json:"email"`
	Password string `gorm:"colmun:password" json:"-"`

	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

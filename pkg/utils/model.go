package utils

import (
	"fmt"

	"gorm.io/gorm"
)

type Criteria interface {
	ToORM(db *gorm.DB) (*gorm.DB, error)
}

type Page struct {
	Page, PageSize int
}

func (c *Page) ToORM(db *gorm.DB) (*gorm.DB, error) {
	if c.Page < 1 {
		c.Page = 1
	}
	if c.PageSize < 1 {
		c.PageSize = 10
	}

	offset := (c.Page - 1) * c.PageSize

	return db.Offset(offset).Limit(c.PageSize), nil
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return func(db *gorm.DB) *gorm.DB {
		crit := &Page{
			Page:     page,
			PageSize: pageSize,
		}
		res, _ := crit.ToORM(db)
		return res
	}
}

type NumberTypes interface {
	~uint | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64
}

func TruncateArray[T NumberTypes](s []T, size uint64) []T {
	return s[:size]
}

func IdsIntToInCon[T NumberTypes](ids []T) string {
	inCon := ""
	for k, v := range ids {
		if k == 0 {
			inCon += fmt.Sprintf("%v", v)
		} else {
			inCon += fmt.Sprintf(",%v", v)
		}
	}
	if inCon == "" {
		inCon = "0"
	}
	return inCon
}

type SIUTypes interface {
	~string | int | uint | int8
}

func Contains[T SIUTypes](s []T, str T) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

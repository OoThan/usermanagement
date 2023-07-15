package utils

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

func IsDuplicate(err error) bool {
	return strings.Contains(err.Error(), "Error 1062") || strings.Contains(err.Error(), "Duplicate entry")
}

func IsErrorNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

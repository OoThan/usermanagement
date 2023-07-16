package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func IsDuplicate(err error) bool {
	return strings.Contains(err.Error(), "Error 1062") || strings.Contains(err.Error(), "Duplicate entry")
}

func IsErrorNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}

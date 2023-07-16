package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/OoThan/usermanagement/pkg/dto"
	"github.com/go-playground/validator/v10"
)

func GenerateSuccessResponse(data any) *dto.Response {
	res := &dto.Response{}
	res.ErrCode = 0
	res.ErrMsg = "Success"
	res.Data = data
	res.HttpStatusCode = http.StatusOK
	return res
}

func GenerateValidationErrorResponse(err error) *dto.Response {
	res := &dto.Response{}
	res.ErrMsg = err.Error()
	if IsValidationError(err) {
		res.ErrCode = 422
		res.ErrMsg = GenerateValidationErrorMessage(err)
		res.HttpStatusCode = http.StatusUnprocessableEntity
		return res
	}

	res.ErrCode = 500
	res.HttpStatusCode = http.StatusInternalServerError
	return res
}

func msgForTag(fe validator.FieldError) string {
	field := CapitalToUnderScore(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%v field is required.", field)
	case "required_with":
		fmt.Println(fe.Param(), "ADSSA")
		return fmt.Sprintf("%v field is required.", field)
	case "oneof":
		return fmt.Sprintf("%v field must be one of %v", field, fe.Param())
	case "email":
		return "Invalid email."
	case "gte", "lte":
		return "invalid length"
	default:
		return "invalid payload" // default error
	}
}

func GenerateValidationErrorMessage(err error) string {
	if vErr, ok := err.(validator.ValidationErrors); ok {
		errMsg := ""
		for _, fieldErr := range vErr {
			errMsg += msgForTag(fieldErr)
		}
		return errMsg
	}

	return err.Error()
}

func GenerateAuthErrorResponse(err error) *dto.Response {
	res := &dto.Response{}
	res.ErrCode = 401
	res.ErrMsg = err.Error()
	res.HttpStatusCode = http.StatusUnauthorized
	return res
}

func GenerateTokenExpireErrorResponse(err error) *dto.Response {
	res := &dto.Response{}
	res.ErrCode = 403
	res.ErrMsg = err.Error()
	res.HttpStatusCode = http.StatusOK
	return res
}

func GenerateGormErrorResponse(err error) *dto.Response {
	res := &dto.Response{}
	res.ErrMsg = err.Error()
	if IsErrorNotFound(err) {
		res.ErrCode = 400
		res.HttpStatusCode = http.StatusBadRequest
		return res
	}

	if IsDuplicate(err) {
		fields := strings.Split(err.Error(), ".")
		field := fields[len(fields)-1]
		res.ErrCode = 400
		msg := "Duplicate Entry"
		res.ErrMsg = fmt.Sprintf("%v %s", strings.Trim(field, "'"), msg)
		res.HttpStatusCode = http.StatusBadRequest
		return res
	}

	res.ErrCode = 500
	res.HttpStatusCode = http.StatusInternalServerError
	return res
}

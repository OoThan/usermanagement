package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/OoThan/usermanagement/pkg/dto"
)

func GenerateSuccessResponse(data any) *dto.Response {
	res := &dto.Response{}
	res.ErrCode = 0
	res.ErrMsg = "Success"
	res.Data = data
	res.HttpStatusCode = http.StatusOK
	return res
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

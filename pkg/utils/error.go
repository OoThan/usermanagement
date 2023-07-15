package utils

import (
	"net/http"

	"github.com/OoThan/usermanagement/pkg/dto"
)

func GenerateAuthErrorResponse(err error) *dto.Response {
	res := &dto.Response{}
	res.ErrCode = 401
	res.ErrMsg = err.Error()
	res.HttpStatusCode = http.StatusUnauthorized
	return res
}

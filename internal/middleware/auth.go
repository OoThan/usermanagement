package middleware

import (
	"fmt"
	"strings"

	"github.com/OoThan/usermanagement/internal/repository"
	"github.com/OoThan/usermanagement/pkg/utils"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	AccessToken string `header:"Authorization"`
}

func AuthMiddleware(r *repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := &authHandler{}
		accessToken := strings.Split(h.AccessToken, "Bearer ")
		if len(accessToken) != 2 {
			res := utils.GenerateAuthErrorResponse(fmt.Errorf("Permission denied"))
			ctx.JSON(res.HttpStatusCode, res)
			ctx.Abort()
			return
		}
	}
}

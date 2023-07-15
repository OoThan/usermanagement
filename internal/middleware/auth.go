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

		claim, err := utils.ValidateAccessToken(accessToken[1])
		if err != nil {
			res := utils.GenerateTokenExpireErrorResponse(fmt.Errorf("token experied"))
			ctx.JSON(res.HttpStatusCode, res)
			ctx.Abort()
			return
		}

		user, err := r.User.FindByField(ctx.Request.Context(), "id", claim.Id)
		if err != nil {
			res := utils.GenerateGormErrorResponse(err)
			ctx.JSON(res.HttpStatusCode, res)
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Set("DS", r.DS.DB)
		ctx.Next()
	}
}

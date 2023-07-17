package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/OoThan/usermanagement/internal/middleware"
	"github.com/OoThan/usermanagement/internal/repository"
	"github.com/OoThan/usermanagement/pkg/dto"
	"github.com/OoThan/usermanagement/pkg/utils"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	R    *gin.Engine
	repo *repository.Repository
}

func newAuthHandler(h *Handler) *authHandler {
	return &authHandler{
		R:    h.R,
		repo: h.repo,
	}
}

func (ctr *authHandler) register() {
	group := ctr.R.Group("/api/auth")
	group.POST("/login", ctr.login)

	group.Use(middleware.AuthMiddleware(ctr.repo))
	group.POST("/logout", ctr.logout)
	group.POST("/refresh", ctr.refresh)
}

func (ctr *authHandler) login(c *gin.Context) {
	req := &dto.UserLoginReq{}
	res := &dto.Response{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	user, err := ctr.repo.User.FindOrByField(c.Request.Context(), "email", "name", req.EmailName)
	if err != nil {
		res = utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	validPass := utils.CheckPasswordHash(req.Password, user.Password)
	if !validPass {
		res.ErrCode = 400
		res.ErrMsg = "Invalid Password"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.Name, user.Id)
	if err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	data := gin.H{
		"access_token": accessToken,
		"name":         user.Name,
	}
	res = utils.GenerateSuccessResponse(data)
	c.JSON(res.HttpStatusCode, res)
}

func (ctr *authHandler) logout(c *gin.Context) {
	c.SetCookie("token", "", 0, "/", c.Request.Host, true, true)

	res := utils.GenerateSuccessResponse(gin.H{
		"message": "logout successful",
	})
	c.JSON(res.HttpStatusCode, res)
}

func (ctr *authHandler) refresh(c *gin.Context) {
	tokens := strings.Split(c.GetHeader("Authorization"), "Bearer ")
	refreshToken, err := utils.GenerateRefreshToken(tokens[1])
	if err != nil {
		// logger.Sugar.Debug(err.Error())
		if strings.Contains(err.Error(), "not expired") {
			res := utils.GenerateSuccessResponse(gin.H{
				"access_token": tokens[1],
			})
			c.JSON(res.HttpStatusCode, res)
			return
		}
		res := utils.GenerateServerError(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	c.SetCookie("token", refreshToken, int(time.Minute)*24, "/", c.Request.Host, true, true)

	res := utils.GenerateSuccessResponse(gin.H{
		"access_token": refreshToken,
	})
	c.JSON(res.HttpStatusCode, res)
}

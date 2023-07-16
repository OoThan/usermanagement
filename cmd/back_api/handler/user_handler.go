package handler

import (
	"github.com/OoThan/usermanagement/internal/repository"
	"github.com/OoThan/usermanagement/pkg/dto"
	"github.com/OoThan/usermanagement/pkg/utils"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	R    *gin.Engine
	repo *repository.Repository
}

func newUserHandler(h *Handler) *userHandler {
	return &userHandler{
		R:    h.R,
		repo: h.repo,
	}
}

func (ctr *userHandler) register() {
	group := ctr.R.Group("/api/users")
	// group.Use(middleware.AuthMiddleware(ctr.repo))

	group.POST("/create", ctr.createUser)
}

func (ctr *userHandler) createUser(c *gin.Context) {
	req := &dto.UserCreateReq{}
	if err := c.ShouldBind(req); err != nil {
		
	}

	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)
}

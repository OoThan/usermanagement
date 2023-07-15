package handler

import (
	"github.com/OoThan/usermanagement/internal/repository"
	"github.com/OoThan/usermanagement/pkg/utils"
	"github.com/gin-gonic/gin"
)

type adminHandler struct {
	R    *gin.Engine
	repo *repository.Repository
}

func newAdminHandler(h *Handler) *adminHandler {
	return &adminHandler{
		R:    h.R,
		repo: h.repo,
	}
}

func (ctr *adminHandler) register() {
	group := ctr.R.Group("/api/users")
	// group.Use(middleware.AuthMiddleware(ctr.repo))

	group.POST("/create", ctr.createAdmin)
}

func (ctr *adminHandler) createAdmin(c *gin.Context) {
	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)
}

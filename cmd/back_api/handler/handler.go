package handler

import (
	"github.com/OoThan/usermanagement/internal/ds"
	"github.com/OoThan/usermanagement/internal/middleware"
	"github.com/OoThan/usermanagement/internal/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	R    *gin.Engine
	repo *repository.Repository
}

type HConfig struct {
	R  *gin.Engine
	DS *ds.DataSource
}

func NewHandler(c *HConfig) *Handler {
	return &Handler{
		R: c.R,
		repo: repository.NewRepository(
			&repository.RepoConfig{
				DS: c.DS,
			},
		),
	}
}

func (h *Handler) Register() {
	h.R.Use(middleware.Cors())

	// auth
	authHandler := newAuthHandler(h)
	authHandler.register()

	// user
	userHandler := newUserHandler(h)
	userHandler.register()
}

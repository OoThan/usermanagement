package handler

import (
	"github.com/OoThan/usermanagement/internal/model"
	"github.com/OoThan/usermanagement/internal/repository"
	"github.com/OoThan/usermanagement/pkg/dto"
	"github.com/OoThan/usermanagement/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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

	
	group.POST("/list", ctr.listUser)
	group.POST("/create", ctr.createUser)
	group.POST("/update", ctr.updateUser)
	group.POST("/delete", ctr.deleteUser)
}

func (ctr *userHandler) createUser(c *gin.Context) {
	req := &dto.UserCreateReq{}
	if err := c.ShouldBind(req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	// hash password
	hashPass, err := utils.HashPassword(req.Password)
	if err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	req.Password = hashPass

	user := &model.User{}
	if err := copier.Copy(&user, req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)
}

func (ctr *userHandler) updateUser(c *gin.Context) {
	req := &dto.UserUpdateReq{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateTokenExpireErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	updateFields := &model.UpdateFields{
		Field: "id",
		Value: req.Id,
		Data:  map[string]any{},
	}

	if req.Password != "" {
		hashPass, err := utils.HashPassword(req.Password)
		if err != nil {
			res := utils.GenerateValidationErrorResponse(err)
			c.JSON(res.HttpStatusCode, res)
			return
		}
		req.Password = hashPass
		updateFields.Data["password"] = hashPass
	}

	updateFields.Data["name"] = req.Name
	updateFields.Data["email"] = req.Email

	if err := ctr.repo.User.Update(c.Request.Context(), updateFields); err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)
}

func (ctr *userHandler) deleteUser(c *gin.Context) {
	req := &dto.UserDeleteReqByIDs{}
	if err := c.ShouldBind(req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	ids := utils.IdsIntToInCon(req.IDS)
	err := ctr.repo.User.Delete(c.Request.Context(), ids)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)
}

func (ctr *userHandler) listUser(c *gin.Context) {
	req := &dto.UserListReq{}
	if err := c.ShouldBind(req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	list, total, err := ctr.repo.User.List(c.Request.Context(), req)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	data := gin.H{
		"list":  list,
		"total": total,
	}
	res := utils.GenerateSuccessResponse(data)
	c.JSON(res.HttpStatusCode, res)
}

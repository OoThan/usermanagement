package handler

import (
	"context"
	"time"

	"github.com/OoThan/usermanagement/internal/middleware"
	"github.com/OoThan/usermanagement/internal/model"
	"github.com/OoThan/usermanagement/internal/repository"
	"github.com/OoThan/usermanagement/pkg/dto"
	"github.com/OoThan/usermanagement/pkg/logger"
	"github.com/OoThan/usermanagement/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
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
	group.Use(middleware.AuthMiddleware(ctr.repo))
	group.POST("/list", ctr.listUser)
	group.POST("/create", ctr.createUser)
	group.POST("/update", ctr.updateUser)
	group.POST("/delete", ctr.deleteUser)
}

func (ctr *userHandler) createUser(c *gin.Context) {
	req := &dto.UserCreateReq{}
	loginUser := c.MustGet("admin").(*model.User)
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

	err = ctr.repo.User.Create(c.Request.Context(), user)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	mdb := ctr.repo.DS.MDB.Database("usermanagement").Collection("user_logs")
	_, err = mdb.InsertOne(context.Background(), bson.D{
		{Key: "user_id", Value: loginUser.Id},
		{Key: "event", Value: "user insert event"},
		{Key: "data", Value: user},
		{Key: "timestamps", Value: time.Now().Unix()},
	})
	if err != nil {
		logger.Sugar.Error(err)
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

	mdb := ctr.repo.DS.MDB.Database("usermanagement").Collection("user_logs")
	_, err := mdb.InsertOne(context.Background(), bson.D{
		{Key: "user_id", Value: updateFields.Value},
		{Key: "event", Value: "user update event"},
		{Key: "data", Value: updateFields.Data},
		{Key: "timestamps", Value: time.Now().Unix()},
	})
	if err != nil {
		logger.Sugar.Error(err)
	}

	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)
}

func (ctr *userHandler) deleteUser(c *gin.Context) {
	req := &dto.UserDeleteReqByIDs{}
	loginUser := c.MustGet("admin").(*model.User)
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

	mdb := ctr.repo.DS.MDB.Database("usermanagement").Collection("user_logs")
	_, err = mdb.InsertOne(context.Background(), bson.D{
		{Key: "user_id", Value: loginUser.Id},
		{Key: "event", Value: "user delete event"},
		{Key: "timestamps", Value: time.Now().Unix()},
	})
	if err != nil {
		logger.Sugar.Error(err)
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

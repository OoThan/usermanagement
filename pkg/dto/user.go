package dto

type UserCreateReq struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserUpdateReq struct {
	Id       uint64 `form:"id" json:"id" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password"`
}

type UserDeleteReqByIDs struct {
	IDS []uint64 `json:"ids" form:"ids" binding:"required,gte=1"`
}

type UserListReq struct {
	ID    uint64 `form:"id" json:"id"`
	Name  string `form:"username" json:"username" `
	Email string `form:"email" json:"email"`
	PageReq
}

type UserLoginReq struct {
	EmailName string `json:"email_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

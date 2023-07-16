package dto

type UserCreateReq struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserListReq struct {
	ID    uint64 `form:"id" json:"id"`
	Name  string `form:"username" json:"username" `
	Email string `form:"email" json:"email"`
	PageReq
}

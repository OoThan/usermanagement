package dto

type UserListReq struct {
	ID    uint64 `form:"id" json:"id"`
	Name  string `form:"username" json:"username" `
	Email string `form:"email" json:"email"`
	PageReq
}

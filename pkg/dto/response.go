package dto

type Response struct {
	ErrCode        uint64 `json:"err_code"`
	ErrMsg         string `json:"err_msg"`
	Data           any    `json:"data,omitempty"`
	HttpStatusCode int    `json:"-"`
}

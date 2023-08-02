package dto

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Avatar   string `json:"avatar"`
}

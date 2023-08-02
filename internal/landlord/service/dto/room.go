package dto

type Room struct {
	ID       int32  `json:"id"`
	Password string `json:"password"`
	Title    string `json:"title"`
}

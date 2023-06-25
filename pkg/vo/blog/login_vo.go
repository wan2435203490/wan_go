package blog

type LoginVO struct {
	Account  string `form:"account"`
	Password string `form:"password"`
	IsAdmin  bool   `form:"isAdmin"`
}

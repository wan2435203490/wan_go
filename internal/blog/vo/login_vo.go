package blog

type LoginVO struct {
	Account  string `form:"username"`
	Password string `form:"password"`
	IsAdmin  bool   `form:"isAdmin"`
}

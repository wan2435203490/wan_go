package blog

type CommentExt struct {
	Comment
	UserName string `json:"userName"`
}

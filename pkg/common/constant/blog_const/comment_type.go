package blog_const

type CommentType struct {
	Code string
	Desc string
}

var (
	COMMENT_TYPE_ARTICLE = CommentType{"article", "文章评论"}
	COMMENT_TYPE_MESSAGE = CommentType{"message", "树洞留言"}
	COMMENT_TYPE_LOVE    = CommentType{"love", "表白墙留言"}
)

func ExistsCommentType(code string) bool {
	switch code {
	case COMMENT_TYPE_ARTICLE.Code:
	case COMMENT_TYPE_MESSAGE.Code:
	case COMMENT_TYPE_LOVE.Code:
		return true
	}
	return false
}

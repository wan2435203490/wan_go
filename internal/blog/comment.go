package blog

import (
	"github.com/gin-gonic/gin"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog/db_comment"
	"wan_go/pkg/common/db/mysql/blog/db_common"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

func deleteCommentCache(source int32, typ string) {
	key := blog_const.COMMENT_COUNT_CACHE + utils.Int32ToString(source) + "_" + typ
	cache.Delete(key)
}

func SaveComment(c *gin.Context) {
	var vo blogVO.CommentVO
	if a.BindFailed(&vo) {
		return
	}
	deleteCommentCache(vo.Source, vo.Type)

	a.DoneApiErr(db_comment.SaveComment(&vo))
}

func DeleteComment(c *gin.Context) {
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	db_comment.DeleteByUserId(id)
	a.OK()
}

func GetCommentCount(c *gin.Context) {
	var source int
	if a.IntFailed(&source, "source") {
		return
	}
	var typ string
	if a.StringFailed(&typ, "type") {
		return
	}

	a.OK(db_common.GetCommentCount(int32(source), typ))
}

func ListComment(c *gin.Context) {
	var vo blogVO.BaseRequestVO[*blogVO.CommentVO]
	if a.BindFailed(&vo) {
		return
	}

	db_comment.ListComment(&vo)

	a.OK(&vo)
}

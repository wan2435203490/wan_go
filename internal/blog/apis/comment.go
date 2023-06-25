package apis

import (
	"github.com/gin-gonic/gin"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog/db_comment"
	"wan_go/pkg/common/db/mysql/blog/db_common"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

type CommentApi struct {
	api.Api
}

func deleteCommentCache(source int32, typ string) {
	key := blog_const.COMMENT_COUNT_CACHE + utils.Int32ToString(source) + "_" + typ
	cache.Delete(key)
}

func (a CommentApi) SaveComment(c *gin.Context) {
	a.MakeContext(c)

	var vo blogVO.CommentVO
	if a.BindFailed(&vo) {
		return
	}
	deleteCommentCache(vo.Source, vo.Type)

	a.DoneApiErr(db_comment.SaveComment(&vo))
}

func (a CommentApi) DeleteComment(c *gin.Context) {
	a.MakeContext(c)
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	db_comment.DeleteByUserId(id)
	a.OK()
}

func (a CommentApi) GetCommentCount(c *gin.Context) {
	a.MakeContext(c)
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

func (a CommentApi) ListComment(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.BaseRequestVO[*blogVO.CommentVO]
	if a.BindPageFailed(&vo) {
		return
	}

	db_comment.ListComment(&vo)

	a.OK(&vo)
}

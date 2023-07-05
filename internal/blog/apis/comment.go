package apis

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/service"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/utils"
	"wan_go/sdk/pkg/jwtauth/user"
)

type CommentApi struct {
	api.Api
}

func deleteCommentCache(source int32, typ string) {
	key := blog_const.COMMENT_COUNT_CACHE + utils.Int32ToString(source) + "_" + typ
	cache.Delete(key)
}

func (a CommentApi) InsertComment(c *gin.Context) {
	s := service.Comment{}
	req := dto.SaveCommentReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	userId, userName := user.GetUserInfo(c)
	req.UserId, req.UserName = userId, userName
	if a.IsError(s.Insert(c, &req)) {
		return
	}
	//deleteCommentCache(req.Source, req.Type)
	a.OKMsg(req.GetId(), constant.DBInsertOK)
}

func (a CommentApi) DeleteComment(c *gin.Context) {
	s := service.Comment{}
	req := dto.DelCommentReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	userId := user.GetUserId32(c)
	if a.IsError(s.Delete(&req, userId)) {
		return
	}
	a.OKMsg(req.GetId(), constant.DBDeleteOK)
}

func (a CommentApi) GetCommentCount(c *gin.Context) {
	s := service.Comment{}
	req := dto.CountCommentReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	var count int64
	if a.IsError(s.CountComment(&req, &count)) {
		return
	}

	a.OK(count)
}

func (a CommentApi) ListComment(c *gin.Context) {
	s := service.Comment{}
	req := dto.PageCommentReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[vo.CommentVO]
	page.Set(req.Pagination)

	if a.IsError(s.Page(c, &req, p, &page)) {
		return
	}

	a.OK(&page)
}

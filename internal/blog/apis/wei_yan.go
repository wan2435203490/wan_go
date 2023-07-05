package apis

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/service"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/sdk/pkg/jwtauth/user"
)

type WeiYanApi struct {
	api.Api
}

func (a WeiYanApi) InsertWeiYan(c *gin.Context) {
	s := service.WeiYan{}
	req := dto.SaveWeiYanReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	req.SetCreateBy(user.GetUserId(c))
	req.UserId = user.GetUserId32(c)
	req.IsPublic = true
	req.Type = blog_const.WEIYAN_TYPE_FRIEND

	if a.IsError(s.InsertReq(&req)) {
		return
	}

	a.OKMsg(req.GetId(), constant.DBInsertOK)
}

func (a WeiYanApi) InsertNews(c *gin.Context) {
	s := service.WeiYan{}
	req := dto.SaveWeiYanReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	var exist bool
	userId := user.GetUserId32(c)
	as := service.Article{}
	if a.IsError(as.ExistArticleByUserId(req.Source, userId, &exist)) {
		return
	}
	if a.IsFailed(!exist, "来源不存在！") {
		return
	}
	req.SetCreateBy(user.GetUserId(c))

	req.UserId = user.GetUserId32(c)
	req.IsPublic = true
	req.Type = blog_const.WEIYAN_TYPE_NEWS

	if a.IsError(s.InsertReq(&req)) {
		return
	}

	a.OKMsg(req.GetId(), constant.DBInsertOK)
}

func (a WeiYanApi) DeleteWeiYan(c *gin.Context) {
	s := service.WeiYan{}
	req := dto.DelWeiYanReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	userId := user.GetUserId32(c)
	if a.IsError(s.Delete(&req, userId)) {
		return
	}

	a.OKMsg(req.GetId(), constant.DBDeleteOK)
}

func (a WeiYanApi) ListNews(c *gin.Context) {
	s := service.WeiYan{}
	req := dto.PageNewsReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.WeiYan]
	page.Set(req.Pagination)

	if a.IsError(s.PageNews(&req, p, &page)) {
		return
	}

	a.OK(&page)
}

func (a WeiYanApi) ListWeiYan(c *gin.Context) {
	s := service.WeiYan{}
	req := dto.PageWeiYanReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.WeiYan]
	page.Set(req.Pagination)

	userId := user.GetUserId32(c)
	if userId == 0 {
		userId = int32(blog_const.ADMIN_USER_ID)
	}
	if a.IsError(s.Page(&req, p, &page, userId)) {
		return
	}

	a.OK(&page)
}

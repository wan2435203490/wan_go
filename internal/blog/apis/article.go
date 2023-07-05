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
	"wan_go/sdk/pkg/jwtauth/user"
)

type ArticleApi struct {
	api.Api
}

func deleteArticleCache(c *gin.Context) {
	key := blog_const.USER_ARTICLE_LIST + user.GetUserIdStr(c)
	cache.Delete(key)
}

func (a ArticleApi) InsertArticle(c *gin.Context) {
	s := service.Article{}
	req := dto.SaveArticleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	userId := user.GetUserId32(c)
	req.UserId = userId
	if a.IsError(s.Insert(&req)) {
		return
	}

	deleteArticleCache(c)
	a.OKMsg(req.GetId(), constant.DBInsertOK)
}

func (a ArticleApi) DeleteArticle(c *gin.Context) {
	s := service.Article{}
	req := dto.DelArticleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	if a.IsError(s.Delete(&req)) {
		return
	}

	deleteArticleCache(c)
	a.OKMsg(req.GetId(), constant.DBDeleteOK)
}

func (a ArticleApi) UpdateArticle(c *gin.Context) {
	s := service.Article{}
	req := dto.SaveArticleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	if a.IsError(s.Update(&req)) {
		return
	}

	deleteArticleCache(c)
	a.OKMsg(req.GetId(), constant.DBUpdateOK)
}

// GetArticle
// flag = true 查询可见的文章
func (a ArticleApi) GetArticle(c *gin.Context) {
	s := service.Article{}
	req := dto.GetArticleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	var data vo.ArticleVO
	userName := user.GetUserName(c)
	if a.IsError(s.GetArticle(c, &req, &data, userName)) {
		return
	}

	a.OKMsg(data, constant.DBGetOK)
}

// ListArticle 查询文章List
func (a ArticleApi) ListArticle(c *gin.Context) {
	s := service.Article{}
	req := dto.PageArticleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	p := actions.GetPermissionFromContext(c)
	var page vo.Page[vo.ArticleVO]
	page.Set(req.Pagination)

	userName := user.GetUserName(c)
	if a.IsError(s.ListArticle(c, &req, p, &page, userName)) {
		return
	}

	a.OK(&page)
}

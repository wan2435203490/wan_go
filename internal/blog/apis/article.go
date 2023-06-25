package apis

import (
	"github.com/gin-gonic/gin"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog/db_article"
	blogVO "wan_go/pkg/vo/blog"
)

type ArticleApi struct {
	api.Api
}

func deleteArticleCache() {
	key := blog_const.USER_ARTICLE_LIST + cache.GetUserIdStr()
	cache.Delete(key)
}

func (a ArticleApi) SaveArticle(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.ArticleVO
	if a.BindFailed(&vo) {
		return
	}
	deleteArticleCache()

	a.Done(db_article.InsertVO(&vo))
}

func (a ArticleApi) DeleteArticle(c *gin.Context) {
	a.MakeContext(c)
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	deleteArticleCache()
	a.Done(db_article.DeleteByUserId(id))
}

func (a ArticleApi) UpdateArticle(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.ArticleVO
	if a.BindFailed(&vo) {
		return
	}

	a.Done(db_article.UpdateVO(&vo))
}

// GetArticleById
// flag = true 查询可见的文章
func (a ArticleApi) GetArticleById(c *gin.Context) {
	a.MakeContext(c)
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	var flag bool
	if a.BoolFailed(&flag, "flag") {
		return
	}
	var password string
	_ = a.StringFailed(&password, "password")

	a.OK(db_article.GetArticleById(id, flag, password))
}

// ListArticle 查询文章List
func (a ArticleApi) ListArticle(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.BaseRequestVO[*blogVO.ArticleVO]
	if a.BindPageFailed(&vo) {
		return
	}

	db_article.ListArticle(&vo)

	a.OK(&vo)
}

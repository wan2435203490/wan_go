package blog

import (
	"github.com/gin-gonic/gin"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog/db_article"
	blogVO "wan_go/pkg/vo/blog"
)

func deleteArticleCache() {
	key := blog_const.USER_ARTICLE_LIST + cache.GetUserIdStr()
	cache.Delete(key)
}

func SaveArticle(c *gin.Context) {
	var vo blogVO.ArticleVO
	if a.BindFailed(&vo) {
		return
	}
	deleteArticleCache()

	a.Done(db_article.InsertVO(&vo))
}

func DeleteArticle(c *gin.Context) {
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	deleteArticleCache()
	a.Done(db_article.DeleteByUserId(id))
}

func UpdateArticle(c *gin.Context) {
	var vo blogVO.ArticleVO
	if a.BindFailed(&vo) {
		return
	}

	a.Done(db_article.UpdateVO(&vo))
}

// GetArticleById
// flag = true 查询可见的文章
func GetArticleById(c *gin.Context) {
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
func ListArticle(c *gin.Context) {
	var vo blogVO.BaseRequestVO[*blogVO.ArticleVO]
	if a.BindFailed(&vo) {
		return
	}

	db_article.ListArticle(&vo)

	a.OK(&vo)
}

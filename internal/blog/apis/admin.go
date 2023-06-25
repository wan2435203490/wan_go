package apis

import (
	"github.com/gin-gonic/gin"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_article"
	"wan_go/pkg/common/db/mysql/blog/db_comment"
	"wan_go/pkg/common/db/mysql/blog/db_tree_hole"
	"wan_go/pkg/common/db/mysql/blog/db_user"
	"wan_go/pkg/common/db/mysql/blog/db_web_info"
	r "wan_go/pkg/common/response"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

type AdminApi struct {
	api.Api
}

func (a AdminApi) ListUser(c *gin.Context) {
	a.MakeContext(c)

	var vo blogVO.BaseRequestVO[*blog.User]
	if a.BindPageFailed(&vo) {
		return
	}

	db_user.ListUser(&vo)

	a.OK(&vo)
}

// ChangeUserStatus
// 修改用户状态
// flag = true：解禁
// flag = false：封禁
func (a AdminApi) ChangeUserStatus(c *gin.Context) {
	a.MakeContext(c)

	var userId int
	if a.IntFailed(&userId, "userId") {
		return
	}
	var userStatus bool
	if a.BoolFailed(&userStatus, "flag") {
		return
	}

	db_user.UpdateUserStatus(userId, userStatus)

	logout(userId)

	a.OK()
}

func logout(userId int) {
	var userIdStr = utils.IntToString(userId)
	deleteTokenCache(blog_const.ADMIN_TOKEN + userIdStr)
	deleteTokenCache(blog_const.USER_TOKEN + userIdStr)
}

func deleteTokenCache(key string) {
	if get, b := cache.GetString(key); b {
		cache.Delete(key)
		cache.Delete(get)
	}
}

// ChangeUserAdmire 修改用户赞赏
func (a AdminApi) ChangeUserAdmire(c *gin.Context) {
	a.MakeContext(c)
	var userId int
	if a.IntFailed(&userId, "userId") {
		return
	}
	var admire string
	if a.StringFailed(&admire, "admire") {
		return
	}

	db_user.UpdateAdmire(userId, admire)

	a.OK()
}

// ChangeUserType 修改用户类型
func (a AdminApi) ChangeUserType(c *gin.Context) {
	a.MakeContext(c)
	var userId int
	if a.IntFailed(&userId, "userId") {
		return
	}
	var userType int
	if a.IntFailed(&userType, "userType") {
		return
	}

	if userType < 0 || userType > 2 {
		a.CodeError(r.PARAMETER_ERROR)
		return
	}

	db_user.UpdateUserType(userId, userType)

	logout(userId)

	a.OK()
}

// GetAdminWebInfo 获取网站信息
func (a AdminApi) GetAdminWebInfo(c *gin.Context) {
	a.MakeContext(c)
	list, err := db_web_info.List()
	if err != nil {
		a.ErrorInternal(err.Error())
		return
	}

	a.OK(list[0])
}

// ListUserArticle 用户查询文章
func (a AdminApi) ListUserArticle(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.BaseRequestVO[*blogVO.ArticleVO]
	if a.BindPageFailed(&vo) {
		return
	}

	db_article.ListAdminArticle(&vo, false)

	a.OK(&vo)
}

// ListBossArticle Boss查询文章
func (a AdminApi) ListBossArticle(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.BaseRequestVO[*blogVO.ArticleVO]
	if a.BindPageFailed(&vo) {
		return
	}

	db_article.ListAdminArticle(&vo, true)

	a.OK(&vo)
}

// ChangeArticleStatus Boss查询文章
func (a AdminApi) ChangeArticleStatus(c *gin.Context) {
	a.MakeContext(c)
	var articleId int
	if a.IntFailed(&articleId, "articleId") {
		return
	}

	var viewStatus, commentStatus, recommendStatus bool
	viewStatusExist := a.BoolFailed(&viewStatus, "viewStatus")
	commentStatusExist := a.BoolFailed(&commentStatus, "commentStatus")
	recommendStatusExist := a.BoolFailed(&recommendStatus, "recommendStatus")

	db_article.ChangeArticleStatus(articleId, viewStatus, commentStatus, recommendStatus,
		viewStatusExist, commentStatusExist, recommendStatusExist)
	a.OK()
}

// GetArticleByIdForUser 查询文章
func (a AdminApi) GetArticleByIdForUser(c *gin.Context) {
	a.MakeContext(c)
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	a.OK(db_article.GetArticleByIdForUser(id))
}

// UserDeleteComment 作者删除评论
func (a AdminApi) UserDeleteComment(c *gin.Context) {
	a.MakeContext(c)
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	a.DoneApiErr(db_comment.UserDeleteComment(id))
}

// BossDeleteComment Boss删除评论
func (a AdminApi) BossDeleteComment(c *gin.Context) {
	a.MakeContext(c)
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	db_comment.Delete(id)
	a.OK()
}

// ListUserComment 用户查询评论
func (a AdminApi) ListUserComment(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.BaseRequestVO[*blog.Comment]
	if a.BindPageFailed(&vo) {
		return
	}

	db_comment.ListAdminComment(&vo, false)

	a.OK(&vo)
}

// ListBossComment Boss查询评论
func (a AdminApi) ListBossComment(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.BaseRequestVO[*blog.Comment]
	if a.BindPageFailed(&vo) {
		return
	}

	db_comment.ListAdminComment(&vo, true)

	a.OK(&vo)
}

// ListBossTreeHole Boss查询树洞
func (a AdminApi) ListBossTreeHole(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.BaseRequestVO[*blog.TreeHole]
	if a.BindPageFailed(&vo) {
		return
	}

	db_tree_hole.ListBossTreeHole(&vo)

	a.OK(&vo)
}

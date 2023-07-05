package apis

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/service"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/db/mysql/blog"
	rocksCache "wan_go/pkg/common/db/rocks_cache"
	r "wan_go/pkg/common/response"
	"wan_go/sdk/pkg/jwtauth/user"
)

type AdminApi struct {
	api.Api
}

// 太多了吧，，，，，，，，，，，，，，，

// ListUser 查询角色列表
//
//	@Summary		查询角色列表
//	@Description	获取JSON
//	@Tags			admin
//	@Param			account			query		string				false	"账号"
//	@Param			userStatus		query		*bool				false	"用户状态"
//	@Param			loginLocation	query		string				false	"归属地"
//	@Param			status			query		string				false	"状态"
//	@Param			userType		query		int					false	"用户类型"
//	@Success		200				{object}	response.Response	"{"code": 200, "data": [...]}"
//	@Router			/admin/user/list [get]
//	@Security		WanBlog
func (a AdminApi) ListUser(c *gin.Context) {
	s := service.Admin{}
	req := dto.PageAdminReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.User]
	page.Set(req.Pagination)

	if a.IsError(s.Page(&req, p, &page)) {
		return
	}

	a.OK(&page)
}

// ChangeUserStatus 修改用户状态
func (a AdminApi) ChangeUserStatus(c *gin.Context) {

	s := service.Admin{}
	req := dto.SaveUserReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if a.IsError(s.UpdateUserStatus(req.UserId, req.UserStatus)) {
		return
	}

	if req.UserStatus == false {
		user.RemoveToken(c)
	}

	a.OK(req.GetId())
}

// ChangeUserAdmire 修改用户赞赏
func (a AdminApi) ChangeUserAdmire(c *gin.Context) {
	s := service.Admin{}
	req := dto.SaveUserReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if a.IsError(s.UpdateUserAdmire(req.UserId, req.Admire)) {
		return
	}

	err := rocksCache.DeleteAdmire()
	if err != nil {
		a.Logger.Errorf("DeleteAdmire error: %s", err)
	}

	a.OK(req.GetId())
}

// ChangeUserType 修改用户类型
func (a AdminApi) ChangeUserType(c *gin.Context) {
	s := service.Admin{}
	req := dto.SaveUserReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if req.UserType < 0 || req.UserType > 2 {
		a.CodeError(r.PARAMETER_ERROR)
		return
	}

	if a.IsError(s.UpdateUserType(req.UserId, req.UserType)) {
		return
	}
	err := rocksCache.DeleteAdmire()
	if err != nil {
		a.Logger.Errorf("DeleteAdmire error: %s", err)
	}

	user.RemoveToken(c)

	a.OK(req.GetId())
}

// GetAdminWebInfo 获取网站信息
func (a AdminApi) GetAdminWebInfo(c *gin.Context) {
	//s := service.Admin{}
	if a.MakeContextChain(c, nil, nil) == nil {
		return
	}

	webInfo, err := rocksCache.GetWebInfo()
	if err != nil {
		a.Logger.Errorf("GetAdminWebInfo error: %s", err)
		return
	}

	a.OK(webInfo[0])
}

// ListArticle 查询文章
func (a AdminApi) ListArticle(c *gin.Context) {
	s := service.Article{}
	req := dto.PageArticleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[vo.ArticleVO]
	page.Set(req.Pagination)

	userId, userName := user.GetUserInfo(c)
	if a.IsError(s.ListAdminArticle(c, &req, p, &page, user.IsAdmin(c), userId, userName)) {
		return
	}

	a.OK(&page)
}

// ChangeArticleStatus Boss查询文章
func (a AdminApi) ChangeArticleStatus(c *gin.Context) {
	s := service.Article{}
	req := dto.ChangeArticleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	userId := user.GetUserId32(c)
	if a.IsError(s.ChangeArticleStatus(&req, userId)) {
		return
	}

	a.OK(req.GetId())
}

// GetArticleByIdForUser 查询文章
func (a AdminApi) GetArticleByIdForUser(c *gin.Context) {
	s := service.Article{}
	req := dto.ChangeArticleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	userId := user.GetUserId32(c)
	data := vo.ArticleVO{}
	if a.IsError(s.GetArticleByIdForUser(req.ArticleId, userId, &data)) {
		return
	}

	a.OK(data)
}

// DeleteComment 删除评论
func (a AdminApi) DeleteComment(c *gin.Context) {
	s := service.Comment{}
	req := dto.DelCommentReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	var err error
	if a.IsAdmin() {
		err = s.BossDeleteComment(&req)
	} else {
		userId := user.GetUserId32(c)
		err = s.UserDeleteComment(&req, userId)
	}
	if a.IsError(err) {
		return
	}

	a.OKMsg(req.GetId(), constant.DBDeleteOK)
}

// ListComment 查询评论
func (a AdminApi) ListComment(c *gin.Context) {
	s := service.Comment{}
	req := dto.PageAdminCommentReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.Comment]
	page.Set(req.Pagination)

	if a.IsError(s.PageAdmin(c, &req, p, &page, user.IsAdmin(c))) {
		return
	}

	a.OK(&page)
}

// ListBossTreeHole Boss查询树洞
func (a AdminApi) ListBossTreeHole(c *gin.Context) {
	s := service.TreeHole{}
	req := dto.PageTreeHoleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.TreeHole]
	page.Set(req.Pagination)

	if a.IsError(s.Page(&req, p, &page)) {
		return
	}

	a.OK(&page)
}

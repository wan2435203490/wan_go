package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/service"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/constant/blog_const"
	rocksCache "wan_go/pkg/common/db/rocks_cache"
	r "wan_go/pkg/common/response"
	"wan_go/pkg/utils"
	"wan_go/sdk/pkg/jwtauth/user"

	"wan_go/pkg/common/db/mysql/blog"
)

type WebInfoApi struct {
	api.Api
}

// UpdateWebInfo 更新网站信息
func (a WebInfoApi) UpdateWebInfo(c *gin.Context) {

	s := service.WebInfo{}
	req := dto.SaveWebInfoReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if a.IsError(s.Update(&req)) {
		return
	}

	if err := rocksCache.DeleteWebInfo(); err != nil {
		a.Logger.Errorf("rocksCache.DeleteWebInfo ID: %d, error:%s", req.GetId(), err.Error())
	}

	a.OK()
}

// GetWebInfo 获取网站信息
func (a WebInfoApi) GetWebInfo(c *gin.Context) {
	s := service.WebInfo{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	webInfos, err := rocksCache.GetWebInfo()
	if err != nil {
		a.Logger.Errorf("GetWebInfo error: %s", err)
	}
	a.OK(webInfos[0])
}

// GetAdmire 获取赞赏
func (a WebInfoApi) GetAdmire(c *gin.Context) {
	s := service.WebInfo{}
	a.MakeContextChain(c, &s.Service, nil)

	//var users []blog.User
	admire, err := rocksCache.GetAdmire()
	if a.IsError(err) {
		return
	}
	a.OK(admire)
}

// GetSortInfo 获取分类标签信息
func (a WebInfoApi) GetSortInfo(c *gin.Context) {
	a.MakeContext(c)
	sortInfos, err := rocksCache.GetSortInfo()
	if err != nil || sortInfos == nil {
		a.ErrorInternal("sortInfos is empty")
		return
	}
	a.OK(sortInfos)
}

// GetWaifuJson 获取看板娘消息
func (a WebInfoApi) GetWaifuJson(c *gin.Context) {
	a.MakeContext(c)
	webInfo, err := rocksCache.GetWebInfo()
	if err != nil || webInfo == nil {
		a.ErrorInternal("webInfo is empty")
		return
	}
	a.Context.String(200, webInfo[0].WaifuJson)
}

// InsertResourcePath 保存
func (a WebInfoApi) InsertResourcePath(c *gin.Context) {

	req := dto.SaveResourcePathReq{}
	s := service.ResourcePath{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if blog_const.RESOURCE_PATH_TYPE_LOVE_PHOTO == req.Type {
		req.Remark = fmt.Sprintf("%d", blog_const.ADMIN_USER_ID)
	}

	if a.IsError(s.Insert(&req)) {
		return
	}

	a.OK(req.GetId(), constant.DBInsertOK)
}

// SaveFriend 保存友链
func (a WebInfoApi) SaveFriend(c *gin.Context) {

	req := dto.SaveResourcePathReq{}
	s := service.ResourcePath{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if utils.IsEmpty(req.Title, req.Cover, req.Url, req.Introduction) {
		a.ErrorInternal("信息不全！")
		return
	}
	req.Type = blog_const.RESOURCE_PATH_TYPE_FRIEND

	if a.IsError(s.Insert(&req)) {
		return
	}

	a.OK(req.GetId(), constant.DBInsertOK)
}

func (a WebInfoApi) DeleteResourcePath(c *gin.Context) {
	req := dto.DelResourcePathReq{}
	s := service.ResourcePath{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if a.IsError(s.Delete(&req)) {
		return
	}

	a.OK(req.GetId(), constant.DBDeleteOK)
}

// UpdateResourcePath 更新
func (a WebInfoApi) UpdateResourcePath(c *gin.Context) {
	req := dto.SaveResourcePathReq{}
	s := service.ResourcePath{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if req.ID == 0 {
		a.ErrorInternal("ID 不能为空")
		return
	}

	if blog_const.RESOURCE_PATH_TYPE_LOVE_PHOTO == req.Type {
		req.Remark = fmt.Sprintf("%d", blog_const.ADMIN_USER_ID)
	}

	if a.IsError(s.Update(&req)) {
		return
	}

	a.OK(req.GetId(), constant.DBUpdateOK)
}

// ListResourcePath 查询资源
func (a WebInfoApi) ListResourcePath(c *gin.Context) {
	req := dto.PageResourcePathReq{}
	s := service.ResourcePath{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if req.Status == nil {
		isAdmin := a.IsAdmin()
		req.Status = &isAdmin
	}
	p := actions.GetPermissionFromContext(c)
	var page vo.Page[vo.ResourcePathVO]
	page.Set(req.Pagination)

	if a.IsError(s.ListByResourceTypeAndClassify(&req, p, &page)) {
		return
	}

	a.OK(page)
}

// ListFunny 查询音乐
//func (a WebInfoApi) ListFunny(c *gin.Context) {
//	a.MakeContext(c)
//	ret := db_resource_path.ListFunny()
//	a.OK(ret)
//}

// ListCollect 查询收藏
func (a WebInfoApi) ListCollect(c *gin.Context) {
	s := service.ResourcePath{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}
	classifyMap := make(map[string][]vo.ResourcePathVO, 0)
	if a.IsError(s.ListCollect(&classifyMap)) {
		return
	}
	a.OK(classifyMap)
}

// SaveFunny 保存音乐
func (a WebInfoApi) SaveFunny(c *gin.Context) {
	req := dto.SaveResourcePathReq{}
	s := service.ResourcePath{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if utils.IsEmpty(req.Classify, req.Cover, req.Url, req.Title) {
		a.ErrorInternal("信息不全！")
		return
	}
	req.Type = blog_const.RESOURCE_PATH_TYPE_FUNNY

	if a.IsError(s.Insert(&req)) {
		return
	}

	a.OK(req.GetId(), constant.DBInsertOK)
}

// ListAdminLovePhoto 查询爱情
func (a WebInfoApi) ListAdminLovePhoto(c *gin.Context) {

	s := service.ResourcePath{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}
	var ret []map[string]any
	if a.IsError(s.ListAdminLovePhoto(blog_const.ADMIN_USER_ID, &ret)) {
		return
	}
	a.OK(&ret)
}

// SaveLovePhoto 保存爱情
func (a WebInfoApi) SaveLovePhoto(c *gin.Context) {
	req := dto.SaveResourcePathReq{}
	s := service.ResourcePath{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if utils.IsEmpty(req.Classify, req.Cover) {
		a.ErrorInternal("信息不全！")
		return
	}
	req.Type = blog_const.RESOURCE_PATH_TYPE_LOVE_PHOTO

	if a.IsError(s.Insert(&req)) {
		return
	}

	a.OK(req.GetId(), constant.DBInsertOK)
}

// InsertTreeHole 树洞
func (a WebInfoApi) InsertTreeHole(c *gin.Context) {
	s := service.TreeHole{}
	req := dto.SaveTreeHoleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if utils.IsEmpty(req.Avatar) {
		req.Avatar = utils.RandomAvatar(user.GetUserId32(c))
	}

	req.SetCreateBy(user.GetUserId(c))
	if a.IsError(s.Insert(&req)) {
		return
	}

	a.OKMsg(req.GetId(), constant.DBInsertOK)
}

func (a WebInfoApi) DeleteTreeHole(c *gin.Context) {
	s := service.TreeHole{}
	req := dto.DelTreeHoleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if a.IsError(s.Delete(&req)) {
		return
	}

	a.OKMsg(req.GetId(), constant.DBDeleteOK)
}

func (a WebInfoApi) ListTreeHole(c *gin.Context) {
	s := service.TreeHole{}
	req := dto.PageTreeHoleReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if req.Pagination == nil {
		req.Pagination = &r.Pagination{Size: 999}
	}
	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.TreeHole]
	page.Set(req.Pagination)

	if a.IsError(s.PageUser(&req, p, &page)) {
		return
	}

	a.OK(&page)
}

// InsertSort 分类
func (a WebInfoApi) InsertSort(c *gin.Context) {
	s := service.Sort{}
	req := dto.SaveSortReq{}

	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	req.SetCreateBy(user.GetUserId(c))
	if a.IsError(s.Insert(&req)) {
		return
	}

	if err := rocksCache.DeleteSortInfo(); err != nil {
		a.Logger.Errorf("rocksCache.DeleteSortInfo ID: %d, error:%s", req.GetId(), err.Error())
	}

	a.OKMsg(req.GetId(), constant.DBInsertOK)
}

func (a WebInfoApi) DeleteSort(c *gin.Context) {
	s := service.Sort{}
	req := dto.DelSortReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if a.IsError(s.Delete(&req)) {
		return
	}

	if err := rocksCache.DeleteSortInfo(); err != nil {
		a.Logger.Errorf("rocksCache.DeleteSortInfo ID: %d, error:%s", req.GetId(), err.Error())
	}

	a.OKMsg(req.GetId(), constant.DBDeleteOK)
}

func (a WebInfoApi) UpdateSort(c *gin.Context) {
	s := service.Sort{}
	req := dto.SaveSortReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if req.ID <= 0 {
		a.ErrorInternal("分类Id不能为空")
		return
	}

	if a.IsError(s.Update(&req)) {
		return
	}

	if err := rocksCache.DeleteSortInfo(); err != nil {
		a.Logger.Errorf("rocksCache.DeleteSortInfo ID: %d, error:%s", req.GetId(), err.Error())
	}

	a.OKMsg(req.GetId(), constant.DBUpdateOK)
}

func (a WebInfoApi) ListSort(c *gin.Context) {
	s := service.Sort{}
	req := dto.PageSortReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.Sort]
	page.Set(req.Pagination)

	if a.IsError(s.Page(&req, p, &page)) {
		return
	}

	a.OK(&page)
}

// InsertLabel 标签
func (a WebInfoApi) InsertLabel(c *gin.Context) {
	s := service.Label{}
	req := dto.SaveLabelReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	req.SetCreateBy(user.GetUserId(c))
	if a.IsError(s.Insert(&req)) {
		return
	}

	if err := rocksCache.DeleteSortInfo(); err != nil {
		a.Logger.Errorf("rocksCache.DeleteSortInfo ID: %d, error:%s", req.GetId(), err.Error())
	}

	a.OKMsg(req.GetId(), constant.DBInsertOK)
}

func (a WebInfoApi) DeleteLabel(c *gin.Context) {
	s := service.Label{}
	req := dto.DelLabelReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if a.IsError(s.Delete(&req)) {
		return
	}

	if err := rocksCache.DeleteSortInfo(); err != nil {
		a.Logger.Errorf("rocksCache.DeleteSortInfo ID: %d, error:%s", req.GetId(), err.Error())
	}

	a.OKMsg(req.GetId(), constant.DBDeleteOK)
}

func (a WebInfoApi) UpdateLabel(c *gin.Context) {
	s := service.Label{}
	req := dto.SaveLabelReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if req.ID <= 0 {
		a.ErrorInternal("标签Id不能为空")
		return
	}

	if a.IsError(s.Update(&req)) {
		return
	}

	if err := rocksCache.DeleteSortInfo(); err != nil {
		a.Logger.Errorf("rocksCache.DeleteSortInfo ID: %d, error:%s", req.GetId(), err.Error())
	}

	a.OKMsg(req.GetId(), constant.DBUpdateOK)
}

func (a WebInfoApi) ListLabel(c *gin.Context) {
	s := service.Label{}
	req := dto.PageLabelReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.Label]
	page.Set(req.Pagination)

	if a.IsError(s.Page(&req, p, &page)) {
		return
	}

	a.OK(&page)
}

func (a WebInfoApi) ListSortAndLabel(c *gin.Context) {
	s := service.Label{}
	req := dto.PageLabelReq{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}
	req.Pagination = &r.Pagination{}
	req.Size = constant.DefaultPageSize

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.Label]
	if a.IsError(s.Page(&req, p, &page)) {
		return
	}

	s2 := service.Sort{}
	req2 := dto.PageSortReq{}
	if a.MakeContextChain(c, &s2.Service, nil) == nil {
		return
	}
	req2.Pagination = &r.Pagination{}
	req2.Size = constant.DefaultPageSize

	var page2 vo.Page[blog.Sort]
	if a.IsError(s2.Page(&req2, p, &page2)) {
		return
	}

	result := make(map[string]any, 2)

	result["labels"] = page.Records
	result["sorts"] = page2.Records

	a.OK(result)
}

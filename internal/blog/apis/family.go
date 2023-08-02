package apis

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"wan_go/internal/blog/service"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	rocksCache "wan_go/pkg/common/db/rocks_cache"
	r "wan_go/pkg/common/response"
	"wan_go/sdk/pkg/jwtauth/user"
)

type FamilyApi struct {
	api.Api
}

func (a FamilyApi) InsertFamily(c *gin.Context) {

	s := service.Family{}
	req := dto.SaveFamilyReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	req.Status = false
	family := blog.Family{}
	req.CopyTo(&family)
	userId := user.GetUserId32(c)
	family.UserId = userId
	if a.IsError(s.Insert(&family)) {
		return
	}

	if a.IsAdmin() {
		err := rocksCache.DeleteAdminFamily()
		if err != nil {
			a.Logger.Errorf("rocksCache.DeleteAdminFamily ID: %d, error:%s", req.GetId(), err.Error())
		}
	}

	a.OK(family.ID, constant.DBInsertOK)
}

func (a FamilyApi) UpdateFamily(c *gin.Context) {
	s := service.Family{}
	req := dto.SaveFamilyReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	userId := user.GetUserId32(c)
	oldFamily := blog.Family{}
	if a.IsError(s.GetByUserId(userId, &oldFamily)) {
		return
	}
	if oldFamily.ID == 0 {
		a.ErrorInternal("家庭信息不存在")
		return
	}

	req.ID = oldFamily.ID
	var family blog.Family
	req.CopyTo(&family)
	if a.IsError(s.Update(&family)) {
		return
	}

	if user.IsAdmin(c) {
		err := rocksCache.DeleteAdminFamily()
		if err != nil {
			a.Logger.Errorf("rocksCache.DeleteAdminFamily ID: %d, error:%s", req.GetId(), err.Error())
		}
	}

	a.OK(family.ID, constant.DBUpdateOK)
}

func (a FamilyApi) DeleteFamily(c *gin.Context) {
	s := service.Family{}
	req := dto.DelFamilyReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	if a.IsError(s.Delete(&req)) {
		return
	}

	a.OKMsg(req.GetId(), constant.DBDeleteOK)
}

func (a FamilyApi) GetFamily(c *gin.Context) {
	s := service.Family{}
	req := dto.DelFamilyReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	userId := user.GetUserId32(c)

	family := blog.Family{}
	if a.IsError(s.GetByUserId(userId, &family)) {
		return
	}

	res := vo.FamilyVO{}
	res.Copy(&family)
	a.OK(res)
}

func (a FamilyApi) GetAdminFamily(c *gin.Context) {
	a.MakeContext(c)

	family, err := rocksCache.GetAdminFamily(constant.AdminID)
	if err != nil || family == nil {
		a.ErrorInternal("家庭信息缓存失败")
		return
	}

	familyVO := vo.FamilyVO{}
	familyVO.Copy(family)
	a.OK(&familyVO)
}

func (a FamilyApi) ListRandomFamily(c *gin.Context) {
	s := service.Family{}
	req := dto.PageFamilyReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	p := actions.GetPermissionFromContext(c)
	if req.Pagination == nil {
		req.Pagination = &r.Pagination{Size: 99}
	}
	var page vo.Page[blog.Family]
	page.Set(req.Pagination)
	if req.Size == 0 {
		req.Size = blog_const.FAMILY_COUNT
	}
	page.Size = req.Size
	status := true
	req.Status = &status
	if a.IsError(s.Page(&req, p, &page)) {
		return
	}

	rand.Shuffle(len(page.Records), func(i, j int) {
		page.Records[i], page.Records[j] = page.Records[j], page.Records[i]
	})

	a.OK(&page)
}

func (a FamilyApi) ListFamily(c *gin.Context) {
	s := service.Family{}
	req := dto.PageFamilyReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.Family]
	page.Set(req.Pagination)
	if req.Size == 0 {
		req.Size = blog_const.FAMILY_COUNT
	}
	page.Size = req.Size
	if a.IsError(s.Page(&req, p, &page)) {
		return
	}
	a.OK(&page)
}

func (a FamilyApi) ChangeLoveStatus(c *gin.Context) {
	s := service.Family{}
	req := dto.ChangeFamilyReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	if a.IsError(s.ChangeLoveStatus(&req)) {
		return
	}
	a.OK(req.ID, constant.DBUpdateOK)
}

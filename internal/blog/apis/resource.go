package apis

import (
	"github.com/gin-gonic/gin"
	"strings"
	"wan_go/internal/blog/service"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/utils"
	"wan_go/sdk/pkg/jwtauth/user"
)

type ResourceApi struct {
	api.Api
}

func (a ResourceApi) InsertResource(c *gin.Context) {

	s := service.Resource{}
	req := dto.SaveResourceReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	req.UserId = user.GetUserId32(c)
	if a.IsError(s.Insert(&req)) {
		return
	}

	a.OK(req.ID, constant.DBInsertOK)
}

func (a ResourceApi) DeleteResource(c *gin.Context) {
	s := service.Resource{}
	req := dto.DelResourceReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	req.Path = strings.ReplaceAll(req.Path, config.Config.Qiniu.Url, "")
	if a.IsError(s.Delete(&req)) {
		return
	}

	utils.DeleteQiniuFile(req.Path)

	a.OK(req.Path, constant.DBDeleteOK)
}

func (a ResourceApi) GetResourceInfo(c *gin.Context) {
	s := service.Resource{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}
	var resources []blog.Resource
	if a.IsError(s.GetResourceInfo(&resources)) {
		return
	}
	if len(resources) > 0 {
		resourceMap := make(map[string]int32, 16)
		keys := make([]string, 0)
		for _, resource := range resources {
			path := strings.ReplaceAll(resource.Path, config.Config.Qiniu.Url, "")
			resourceMap[path] = resource.ID
			keys = append(keys, path)
		}
		fileInfo := utils.BatchGetFiles(keys)
		if len(fileInfo) > 0 {
			var collect []blog.Resource
			for k, v := range fileInfo {
				res := blog.Resource{
					ID:       resourceMap[k],
					Size:     utils.StringToInt32(v["size"]),
					MimeType: v["mimeType"],
				}
				collect = append(collect, res)
			}
			if a.IsError(s.BatchUpdate(&collect)) {
				return
			}
		}
	}
	a.OK()
}

func (a ResourceApi) GetImageList(c *gin.Context) {
	s := service.Resource{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	var resources []blog.Resource
	adminId := blog_const.ADMIN_USER_ID
	if a.IsError(s.GetImageList(&resources, adminId)) {
		return
	}

	var paths []string
	if len(resources) > 0 {
		for _, res := range resources {
			paths = append(paths, res.Path)
		}
	}

	a.OK(&paths)
}

func (a ResourceApi) ListResource(c *gin.Context) {
	s := service.Resource{}
	req := dto.PageResourceReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	p := actions.GetPermissionFromContext(c)
	var page vo.Page[blog.Resource]
	page.Set(req.Pagination)

	if a.IsError(s.Page(&req, p, &page)) {
		return
	}

	a.OK(&page)
}

func (a ResourceApi) ChangeResourceStatus(c *gin.Context) {
	s := service.Resource{}
	req := dto.ChangeResourceReq{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}
	if a.IsError(s.ChangeResourceStatus(req.ID, *req.Status)) {
		return
	}
	a.OK(req.ID, constant.DBUpdateOK)
}

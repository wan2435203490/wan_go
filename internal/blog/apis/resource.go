package apis

import (
	"github.com/gin-gonic/gin"
	"strings"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_resource"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

type ResourceApi struct {
	api.Api
}

func (a ResourceApi) SaveResource(c *gin.Context) {
	a.MakeContext(c)
	var vo blog.Resource
	if a.BindFailed(&vo) {
		return
	}

	if a.EmptyFailed("资源类型和资源路径不能为空！", vo.Type, vo.Path) {
		return
	}

	resource := blog.Resource{
		UserId:   int32(cache.GetUserId()),
		Path:     vo.Path,
		Size:     vo.Size,
		Type:     vo.Type,
		MimeType: vo.MimeType,
	}

	if a.IsError(db_resource.Insert(&resource)) {
		return
	}
	a.OK()
}

func (a ResourceApi) DeleteResource(c *gin.Context) {
	a.MakeContext(c)
	var path string
	if a.StringFailed(&path, "path") {
		return
	}

	path = strings.ReplaceAll(path, config.Config.Qiniu.Url, "")

	utils.DeleteQiniuFile(path)

	a.Done(db_resource.DeleteByPath(path))
}

func (a ResourceApi) GetResourceInfo(c *gin.Context) {
	a.MakeContext(c)
	resources := db_resource.GetResourceInfo()
	if resources != nil {
		resourceMap := make(map[string]int32, 16)
		keys := make([]string, 0)
		for _, resource := range *resources {
			path := strings.ReplaceAll(resource.Path, config.Config.Qiniu.Url, "")
			resourceMap[path] = resource.ID
			keys = append(keys, path)
		}
		fileInfo := utils.BatchGetFiles(keys)
		if len(fileInfo) > 0 {
			var collect []*blog.Resource
			for k, v := range fileInfo {
				res := blog.Resource{
					ID:       resourceMap[k],
					Size:     utils.StringToInt32(v["size"]),
					MimeType: v["mimeType"],
				}
				collect = append(collect, &res)
			}
			if a.IsError(db_resource.BatchUpdate(&collect)) {
				return
			}
		}
	}
	a.OK()
}

func (a ResourceApi) GetImageList(c *gin.Context) {
	a.MakeContext(c)

	list := db_resource.GetImageList()

	var paths []string
	if len(*list) > 0 {
		for _, res := range *list {
			paths = append(paths, res.Path)
		}
	}

	a.OK(&paths)
}

func (a ResourceApi) ListResource(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.BaseRequestVO[*blog.Resource]
	if a.BindPageFailed(&vo) {
		return
	}

	db_resource.ListResource(&vo)

	a.OK(&vo)
}

func (a ResourceApi) ChangeResourceStatus(c *gin.Context) {
	a.MakeContext(c)
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	var status bool
	if a.BoolFailed(&status, "flag") {
		return
	}

	db_resource.ChangeResourceStatus(id, status)

	a.OK()
}

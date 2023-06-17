package blog

import (
	"github.com/gin-gonic/gin"
	"strings"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_resource"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

func SaveResource(c *gin.Context) {
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

func DeleteResource(c *gin.Context) {
	var path string
	if a.StringFailed(&path, "path") {
		return
	}

	path = strings.ReplaceAll(path, blog_const.DOWNLOAD_URL, "")

	utils.DeleteQiniuFile(&[]string{path})

	a.Done(db_resource.DeleteByPath(path))
}

func GetResourceInfo(c *gin.Context) {
	resources := db_resource.GetResourceInfo()
	if resources != nil {
		resourceMap := make(map[string]int32, 16)
		keys := make([]string, 0)
		for _, resource := range *resources {
			path := strings.ReplaceAll(resource.Path, blog_const.DOWNLOAD_URL, "")
			resourceMap[path] = resource.ID
			keys = append(keys, path)
		}
		fileInfo := utils.GetQiniuFileInfo(&keys)
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

func GetImageList(c *gin.Context) {

	list := db_resource.GetImageList()

	var paths []string
	if len(*list) > 0 {
		for _, res := range *list {
			paths = append(paths, res.Path)
		}
	}

	a.OK(&paths)
}

func ListResource(c *gin.Context) {
	var vo blogVO.BaseRequestVO[*blog.Resource]
	if a.BindFailed(&vo) {
		return
	}

	db_resource.ListResource(&vo)

	a.OK(&vo)
}

func ChangeResourceStatus(c *gin.Context) {
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

package blog

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog/db_common"
	"wan_go/pkg/common/db/mysql/blog/db_label"
	"wan_go/pkg/common/db/mysql/blog/db_resource_path"
	"wan_go/pkg/common/db/mysql/blog/db_sort"
	"wan_go/pkg/common/db/mysql/blog/db_tree_hole"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"

	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_web_info"
	"wan_go/sdk/api"
)

var (
	a webInfoApi
)

type webInfoApi struct {
	api.Api
}

func WebInfo(c *gin.Context) {
	if a.MakeContext(c) != nil {
		c.Abort()
		return
	}

	c.Next()
}

// UpdateWebInfo 更新网站信息
func UpdateWebInfo(c *gin.Context) {
	var webInfo blog.WebInfo
	if a.BindFailed(&webInfo) {
		return
	}

	if a.IsError(db_web_info.Update(&webInfo)) {
		return
	}

	list, err := db_web_info.List()
	if a.IsError(err) {
		return
	}

	//设置没有过期时间的KEY 带上过期时间用cache.DefaultExpiration
	cache.Set(blog_const.WEB_INFO, list[0])

	a.Success()
}

// GetWebInfo 获取网站信息
func GetWebInfo(c *gin.Context) {

	get, b := cache.Get(blog_const.WEB_INFO)
	if b {
		webInfo := get.(*blog.WebInfo)
		var ret *blog.WebInfo
		ret.Copy(webInfo)

		ret.RandomName = ""
		ret.RandomAvatar = ""
		ret.RandomCover = ""
		ret.WaifuJson = ""

		a.OK(ret)
	}

	a.Success()
}

// GetAdmire 获取赞赏
func GetAdmire(c *gin.Context) {
	admire := db_common.GetAdmire()
	a.OK(admire)
}

// GetSortInfo 获取分类标签信息
func GetSortInfo(c *gin.Context) {
	if get, b := cache.Get(blog_const.SORT_INFO); b {
		a.OK(get.([]*blog.Sort))
	}

	a.Success()
}

// GetWifeJson 获取看板娘消息
func GetWifeJson(c *gin.Context) {
	if get, b := cache.Get(blog_const.WEB_INFO); b {
		info := get.(string)
		if len(strings.TrimSpace(info)) > 0 {
			a.OK(info)
		}
	}

	a.OK("{}")
}

// SaveResourcePath 保存
func SaveResourcePath(c *gin.Context) {
	var resourcePathVO blogVO.ResourcePathVO
	if a.BindFailed(&resourcePathVO) {
		return
	}

	if utils.IsEmpty(resourcePathVO.Type) {
		a.ErrorInternal("资源类型不能为空！")
		return
	}

	if blog_const.RESOURCE_PATH_TYPE_LOVE_PHOTO == resourcePathVO.Type {
		admin := cache.GetAdminUser()
		if admin != nil {
			resourcePathVO.Remark = fmt.Sprintf("%d", admin.ID)
		}
	}

	resourcePath := blog.ResourcePath{}
	resourcePath.Copy(&resourcePathVO)

	if a.IsError(db_resource_path.Insert(&resourcePath)) {
		return
	}

	a.Success()
}

// SaveFriend 保存友链
func SaveFriend(c *gin.Context) {

	var vo blogVO.ResourcePathVO
	if a.BindFailed(&vo) {
		return
	}

	if utils.IsEmpty(vo.Title, vo.Cover, vo.Url, vo.Introduction) {
		a.ErrorInternal("信息不全！")
		return
	}

	path := copyResourcePath(&vo)

	if a.IsError(db_resource_path.Insert(path)) {
		return
	}

	a.Success()
}

func copyResourcePath(vo *blogVO.ResourcePathVO) *blog.ResourcePath {
	path := blog.ResourcePath{}

	path.Title = vo.Title
	path.Introduction = vo.Introduction
	path.Cover = vo.Cover
	path.Url = vo.Url
	path.Type = blog_const.RESOURCE_PATH_TYPE_FRIEND
	path.Status = false

	return &path
}

func DeleteResourcePath(c *gin.Context) {
	id := a.Param("id")

	if id == "" {
		a.ErrorInternal("ResourcePathId is empty")
		return
	}

	path := blog.ResourcePath{ID: utils.StringToInt32(id)}

	a.Done(db_resource_path.Delete(&path))
}

// UpdateResourcePath 更新
func UpdateResourcePath(c *gin.Context) {
	var resourcePathVO blogVO.ResourcePathVO
	if a.BindFailed(&resourcePathVO) {
		return
	}

	if utils.IsEmpty(resourcePathVO.Type) {
		a.ErrorInternal("资源类型不能为空！")
		return
	}
	if resourcePathVO.ID == 0 {
		a.ErrorInternal("ID不能为空！")
	}

	if blog_const.RESOURCE_PATH_TYPE_LOVE_PHOTO == resourcePathVO.Type {
		admin := cache.GetAdminUser()
		if admin != nil {
			resourcePathVO.Remark = fmt.Sprintf("%d", admin.ID)
		}
	}

	resourcePath := blog.ResourcePath{}
	resourcePath.Copy(&resourcePathVO)

	if a.IsError(db_resource_path.Update(&resourcePath)) {
		return
	}

	a.Success()
}

// ListResourcePath 查询资源
func ListResourcePath(c *gin.Context) {
	var requestVO blogVO.BaseRequestVO[*blogVO.ResourcePathVO]
	if a.BindFailed(&requestVO) {
		return
	}

	requestVO.Status = a.IsAdmin() || requestVO.Status

	db_resource_path.ListByResourceTypeAndClassify(&requestVO)

	a.OK(requestVO)
}

// ListFunny 查询音乐
func ListFunny(c *gin.Context) {
	ret := db_resource_path.ListFunny()
	a.OK(ret)
}

// ListCollect 查询收藏
func ListCollect(c *gin.Context) {
	ret := db_resource_path.ListCollect()
	a.OK(ret)
}

// SaveFunny 保存音乐
func SaveFunny(c *gin.Context) {
	var vo blogVO.ResourcePathVO
	if a.BindFailed(&vo) {
		return
	}

	if a.IsFailed(utils.IsEmpty(vo.Classify, vo.Cover, vo.Url, vo.Title), "信息不全！") {
		return
	}

	path := blog.ResourcePath{}
	path.Classify = vo.Classify
	path.Title = vo.Title
	path.Cover = vo.Cover
	path.Url = vo.Url
	path.Type = blog_const.RESOURCE_PATH_TYPE_FUNNY

	a.Done(db_resource_path.Insert(&path))
}

// ListAdminLovePhoto 查询爱情
func ListAdminLovePhoto(c *gin.Context) {
	ret := db_resource_path.ListAdminLovePhoto(a.AdminId())
	a.OK(ret)
}

// SaveLovePhoto 保存爱情
func SaveLovePhoto(c *gin.Context) {
	var vo blogVO.ResourcePathVO
	if a.BindFailed(&vo) {
		return
	}

	if utils.IsEmpty(vo.Classify, vo.Cover) {
		a.ErrorInternal("信息不全！")
	}
	path := blog.ResourcePath{}
	path.Classify = vo.Classify
	path.Title = vo.Title
	path.Cover = vo.Cover
	path.Remark = utils.Int32ToString(a.GetCurrentUserId())
	path.Type = blog_const.RESOURCE_PATH_TYPE_LOVE_PHOTO

	a.Done(db_resource_path.Insert(&path))
}

// SaveTreeHole 树洞
func SaveTreeHole(c *gin.Context) {
	var treeHole blog.TreeHole
	if a.BindFailed(&treeHole) {
		return
	}

	if a.IsFailed(utils.IsEmpty(treeHole.Message), "留言不能为空！") {
		return
	}

	if a.IsError(db_tree_hole.Insert(&treeHole)) {
		return
	}

	if utils.IsEmpty(treeHole.Avatar) {
		//todo
		//treeHole.Avatar = randomAvatars
	}

	a.Success()

}

func DeleteTreeHole(c *gin.Context) {
	id := utils.StringToInt32(a.Param("id"))
	a.Done(db_tree_hole.Delete(&blog.TreeHole{ID: id}))
}

func ListTreeHole(c *gin.Context) {

	treeHoles := db_tree_hole.List()
	if treeHoles == nil {
		a.Success()
	}

	for _, treeHole := range treeHoles {
		if utils.IsEmpty(treeHole.Avatar) {
			//todo
			//treeHole.Avatar = randomAvatars
		}
	}

	a.OK(treeHoles)
}

// SaveSort 分类
func SaveSort(c *gin.Context) {
	var sort blog.Sort
	if a.BindFailed(&sort) {
		return
	}

	if a.EmptyFailed("分类名称和分类描述不能为空！", sort.SortName, sort.SortDescription) {
		return
	}

	if a.IsFailed(sort.SortType == blog_const.SORT_TYPE_BAR.Code && sort.Priority == 0,
		"导航栏分类必须配置优先级！") {
		return
	}

	if a.IsError(db_sort.Insert(&sort)) {
		return
	}

	db_common.CacheSort()

	a.Success()
}

func DeleteSort(c *gin.Context) {
	id := utils.StringToInt32(a.Param("id"))

	err := db_sort.Delete(&blog.Sort{ID: id})

	db_common.CacheSort()

	a.Done(err)
}

func UpdateSort(c *gin.Context) {
	var sort blog.Sort
	if a.BindFailed(&sort) {
		return
	}

	err := db_sort.Update(&sort)

	db_common.CacheSort()

	a.Done(err)
}

func ListSort(c *gin.Context) {
	a.OK(db_sort.List())
}

// SaveLabel 标签
func SaveLabel(c *gin.Context) {
	var label blog.Label
	if a.BindFailed(&label) {
		return
	}

	if a.EmptyFailed("标签名称和标签描述不能为空！", label.LabelName, label.LabelDescription) {
		return
	}

	if a.IsFailed(label.SortId == 0, "分类Id不能为空！") {
		return
	}

	if a.IsError(db_label.Insert(&label)) {
		return
	}

	db_common.CacheSort()

	a.Success()
}

func DeleteLabel(c *gin.Context) {
	id := utils.StringToInt32(a.Param("id"))

	err := db_label.Delete(&blog.Label{ID: id})

	db_common.CacheSort()

	a.Done(err)
}

func UpdateLabel(c *gin.Context) {
	var label blog.Label
	if a.BindFailed(&label) {
		return
	}

	err := db_label.Update(&label)

	db_common.CacheSort()

	a.Done(err)
}

func ListLabel(c *gin.Context) {
	a.OK(db_label.List())
}

func ListSortAndLabel(c *gin.Context) {
	sorts := db_sort.List()
	labels := db_label.List()

	result := make(map[string]any, 2)

	result["sorts"] = sorts
	result["labels"] = labels

	a.OK(result)
}

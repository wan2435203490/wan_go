package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strings"
	"wan_go/pkg/common/api"
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
)

type WebInfoApi struct {
	api.Api
}

// UpdateWebInfo 更新网站信息
func (a WebInfoApi) UpdateWebInfo(c *gin.Context) {
	a.MakeContext(c)
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

	a.OK()
}

// GetWebInfo 获取网站信息
func (a WebInfoApi) GetWebInfo(c *gin.Context) {
	a.MakeContext(c)

	get, b := cache.Get(blog_const.WEB_INFO)
	if b {
		webInfo := get.(*blog.WebInfo)
		var ret blog.WebInfo
		ret.Copy(webInfo)

		ret.RandomName = ""
		ret.RandomAvatar = ""
		ret.RandomCover = ""
		ret.WaifuJson = ""

		a.OK(ret)
		return
	}

	a.OK()
}

// GetAdmire 获取赞赏
func (a WebInfoApi) GetAdmire(c *gin.Context) {
	a.MakeContext(c)
	admire := db_common.GetAdmire()
	a.OK(admire)
}

// GetSortInfo 获取分类标签信息
func (a WebInfoApi) GetSortInfo(c *gin.Context) {
	a.MakeContext(c)
	if get, b := cache.Get(blog_const.SORT_INFO); b {
		a.OK(get.(*[]*blog.Sort))
		return
	}

	a.OK()
}

// GetWifeJson 获取看板娘消息
func (a WebInfoApi) GetWifeJson(c *gin.Context) {
	a.MakeContext(c)
	if get, b := cache.Get(blog_const.WEB_INFO); b {
		info := get.(string)
		if len(strings.TrimSpace(info)) > 0 {
			a.OK(info)
			return
		}
	}

	a.OK("{}")
}

// SaveResourcePath 保存
func (a WebInfoApi) SaveResourcePath(c *gin.Context) {
	a.MakeContext(c)
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
	resourcePathVO.CopyTo(&resourcePath)

	if a.IsError(db_resource_path.Insert(&resourcePath)) {
		return
	}

	a.OK()
}

// SaveFriend 保存友链
func (a WebInfoApi) SaveFriend(c *gin.Context) {
	a.MakeContext(c)

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

	a.OK()
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

func (a WebInfoApi) DeleteResourcePath(c *gin.Context) {
	a.MakeContext(c)
	id := a.Param("id")

	if id == "" {
		a.ErrorInternal("ResourcePathId is empty")
		return
	}

	path := blog.ResourcePath{ID: utils.StringToInt32(id)}

	a.Done(db_resource_path.Delete(&path))
}

// UpdateResourcePath 更新
func (a WebInfoApi) UpdateResourcePath(c *gin.Context) {
	a.MakeContext(c)
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
	resourcePathVO.CopyTo(&resourcePath)

	if a.IsError(db_resource_path.Update(&resourcePath)) {
		return
	}

	a.OK()
}

// ListResourcePath 查询资源
func (a WebInfoApi) ListResourcePath(c *gin.Context) {
	a.MakeContext(c)
	var requestVO blogVO.BaseRequestVO[*blogVO.ResourcePathVO]
	if a.BindPageFailed(&requestVO) {
		return
	}

	requestVO.Status = a.IsAdmin() || requestVO.Status

	db_resource_path.ListByResourceTypeAndClassify(&requestVO)

	a.OK(requestVO)
}

// ListFunny 查询音乐
func (a WebInfoApi) ListFunny(c *gin.Context) {
	a.MakeContext(c)
	ret := db_resource_path.ListFunny()
	a.OK(ret)
}

// ListCollect 查询收藏
func (a WebInfoApi) ListCollect(c *gin.Context) {
	a.MakeContext(c)
	ret := db_resource_path.ListCollect()
	a.OK(ret)
}

// SaveFunny 保存音乐
func (a WebInfoApi) SaveFunny(c *gin.Context) {
	a.MakeContext(c)
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
func (a WebInfoApi) ListAdminLovePhoto(c *gin.Context) {
	a.MakeContext(c)
	ret := db_resource_path.ListAdminLovePhoto(a.AdminId())
	a.OK(ret)
}

// SaveLovePhoto 保存爱情
func (a WebInfoApi) SaveLovePhoto(c *gin.Context) {
	a.MakeContext(c)
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
func (a WebInfoApi) SaveTreeHole(c *gin.Context) {
	a.MakeContext(c)
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
		//treeHole.Avatar = utils.RandomAvatars()
	}

	a.OK()

}

func (a WebInfoApi) DeleteTreeHole(c *gin.Context) {
	a.MakeContext(c)
	id := utils.StringToInt32(a.Param("id"))
	a.Done(db_tree_hole.Delete(&blog.TreeHole{ID: id}))
}

func (a WebInfoApi) ListTreeHole(c *gin.Context) {
	a.MakeContext(c)

	treeHoles := db_tree_hole.List()
	if treeHoles == nil {
		a.OK()
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
func (a WebInfoApi) SaveSort(c *gin.Context) {
	a.MakeContext(c)
	var sort blog.Sort
	if a.BindFailed(&sort, binding.JSON) {
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

	a.OK()
}

func (a WebInfoApi) DeleteSort(c *gin.Context) {
	a.MakeContext(c)
	id := utils.StringToInt32(a.Param("id"))

	err := db_sort.Delete(&blog.Sort{ID: id})

	db_common.CacheSort()

	a.Done(err)
}

func (a WebInfoApi) UpdateSort(c *gin.Context) {
	a.MakeContext(c)
	var sort blog.Sort
	if a.BindFailed(&sort) {
		return
	}

	err := db_sort.Update(&sort)

	db_common.CacheSort()

	a.Done(err)
}

func (a WebInfoApi) ListSort(c *gin.Context) {
	a.MakeContext(c)
	a.OK(db_sort.List())
}

// SaveLabel 标签
func (a WebInfoApi) SaveLabel(c *gin.Context) {
	a.MakeContext(c)
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

	a.OK()
}

func (a WebInfoApi) DeleteLabel(c *gin.Context) {
	a.MakeContext(c)
	id := utils.StringToInt32(a.Param("id"))

	err := db_label.Delete(&blog.Label{ID: id})

	db_common.CacheSort()

	a.Done(err)
}

func (a WebInfoApi) UpdateLabel(c *gin.Context) {
	a.MakeContext(c)
	var label blog.Label
	if a.BindFailed(&label) {
		return
	}

	err := db_label.Update(&label)

	db_common.CacheSort()

	a.Done(err)
}

func (a WebInfoApi) ListLabel(c *gin.Context) {
	a.MakeContext(c)
	a.OK(db_label.List())
}

func (a WebInfoApi) ListSortAndLabel(c *gin.Context) {
	a.MakeContext(c)
	sorts := db_sort.List()
	labels := db_label.List()

	result := make(map[string]any, 2)

	result["sorts"] = sorts
	result["labels"] = labels

	a.OK(result)
}

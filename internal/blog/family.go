package blog

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_common"
	"wan_go/pkg/common/db/mysql/blog/db_family"
	blogVO "wan_go/pkg/vo/blog"
)

func deleteFamilyListCache() {
	//cache.Delete(blog_const.ADMIN_FAMILY)
	cache.Delete(blog_const.FAMILY_LIST)
}

func cacheAdminFamily(family *blog.Family) {
	cache.Set(blog_const.ADMIN_FAMILY, family)
}

func SaveFamily(c *gin.Context) {
	var vo blogVO.FamilyVO
	if a.BindFailed(&vo) {
		return
	}

	userId := cache.GetUserId()
	vo.UserId = int32(userId)

	oldFamily := db_family.GetByUserId(userId)

	family := blog.Family{}
	vo.CopyTo(&family)
	family.Status = false

	if oldFamily == nil {
		family.ID = oldFamily.ID
		db_family.Update(&family)
	} else {
		db_family.Insert(&family)
	}

	if userId == cache.GetAdminUserId() {
		cacheAdminFamily(&family)
	}

	deleteFamilyListCache()

	a.OK()
}

func DeleteFamily(c *gin.Context) {
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	db_family.DeleteById(id)
	deleteFamilyListCache()
	a.OK()
}

func GetFamily(c *gin.Context) {
	userId := cache.GetUserId()
	family := db_family.GetByUserId(userId)

	if family == nil {
		a.OK()
	} else {
		vo := blogVO.FamilyVO{}
		vo.Copy(family)
		a.OK(&vo)
	}
}

func GetAdminFamily(c *gin.Context) {

	family, ok := cache.Get(blog_const.ADMIN_FAMILY)
	if !ok {
		a.ErrorInternal("家庭信息缓存失败")
		return
	}

	if family == nil {
		a.OK()
	} else {
		vo := blogVO.FamilyVO{}
		vo.Copy(family.(*blog.Family))
		a.OK(&vo)
	}
}

func ListRandomFamily(c *gin.Context) {

	var size int
	if a.IntFailed(&size, "size") {
		size = 10
	}

	list := db_common.GetFamilyList()
	if list == nil {
		a.ErrorInternal("找不到家庭信息")
		return
	}
	n := len(list)
	if n > size {
		rand.Shuffle(n, func(i, j int) {
			list[i], list[j] = list[j], list[i]
		})
		list = list[:size]
	}

	a.OK(&list)
}

func ListFamily(c *gin.Context) {
	var vo blogVO.BaseRequestVO[*blog.Family]
	if a.BindFailed(&vo) {
		return
	}

	db_family.ListFamily(&vo)

	a.OK(&vo)
}

func ChangeLoveStatus(c *gin.Context) {
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	var status bool
	if a.BoolFailed(&status, "flag") {
		return
	}
	db_family.ChangeLoveStatus(id, status)
	deleteFamilyListCache()
	a.OK()
}

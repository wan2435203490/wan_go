package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_article"
	"wan_go/pkg/common/db/mysql/blog/db_wei_yan"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

type WeiYanApi struct {
	api.Api
}

func (a WeiYanApi) SaveWeiYan(c *gin.Context) {
	a.MakeContext(c)
	var vo blog.WeiYan
	if a.BindFailed(&vo) {
		return
	}

	if a.EmptyFailed("微言不能为空！", vo.Content) {
		return
	}

	userId := cache.GetUserId()
	if a.IsFailed(userId < 1, "用户登录信息失效，请重新登录！") {
		return
	}

	weiYan := blog.WeiYan{
		UserId:   int32(userId),
		Content:  vo.Content,
		IsPublic: vo.IsPublic,
		Type:     blog_const.WEIYAN_TYPE_FRIEND,
	}

	if a.IsError(db_wei_yan.Insert(&weiYan)) {
		return
	}
	a.OK()
}

func (a WeiYanApi) SaveNews(c *gin.Context) {
	a.MakeContext(c)
	var vo blog.WeiYan
	if a.BindFailed(&vo) {
		return
	}

	if a.IsFailed(utils.IsEmpty(vo.Content) || vo.Source == 0 || vo.CreatedAt.IsZero(), "信息不全！") {
		return
	}

	exist := db_article.ExistArticleByUserId(vo.Source)
	if a.IsFailed(!exist, "来源不存在！") {
		return
	}

	weiYan := blog.WeiYan{
		UserId:   int32(cache.GetUserId()),
		Content:  vo.Content,
		IsPublic: true,
		Source:   vo.Source,
		Type:     blog_const.WEIYAN_TYPE_NEWS,
	}

	if a.IsError(db_wei_yan.Insert(&weiYan)) {
		return
	}
	a.OK()
}

func (a WeiYanApi) DeleteWeiYan(c *gin.Context) {
	a.MakeContext(c)
	var id int
	if a.IntFailed(&id, "id") {
		return
	}
	a.Done(db_wei_yan.DeleteByUserId(id))
}

func (a WeiYanApi) ListNews(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.BaseRequestVO[*blog.WeiYan]
	if a.BindPageFailed(&vo) {
		return
	}

	if a.IsFailed(vo.Source == 0, "来源不能为空！") {
		return
	}

	db_wei_yan.ListNews(&vo)

	a.OK(&vo)
}

func (a WeiYanApi) ListWeiYan(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.BaseRequestVO[*blog.WeiYan]
	if a.BindPageFailed(&vo) {
		return
	}

	db_wei_yan.ListWeiYan(&vo)

	fmt.Println("qwqwq")
	a.OK(&vo)
}

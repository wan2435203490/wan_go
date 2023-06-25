package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
)

func init() {
	routerCheckRole = append(routerCheckRole, weiYanRouter)
}

func weiYanRouter(v1 *gin.RouterGroup) {
	api := apis.WeiYanApi{}
	r := v1.Group("/weiYan")
	{
		r.POST("/saveWeiYan", api.SaveWeiYan)
		r.POST("/saveNews", api.SaveNews)
		r.GET("/listNews", api.ListNews)
		r.GET("/deleteWeiYan", api.DeleteWeiYan)
		r.GET("/listWeiYan", api.ListWeiYan)
	}
}

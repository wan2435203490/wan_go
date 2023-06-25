package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, qiniuRouter)
}

func qiniuRouter(v1 *gin.RouterGroup) {
	api := apis.QiniuApi{}
	r := v1.Group("/qiniu").Use(middleware.LoginCheck)
	{
		r.GET("/getUpToken", api.GetUpToken) //0
		//r.GET("/testUpload", api.TestUpload) //0
	}
}

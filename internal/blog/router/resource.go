package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, resourceRouter)
}

func resourceRouter(v1 *gin.RouterGroup) {
	api := apis.ResourceApi{}
	r := v1.Group("/resource").Use(middleware.LoginCheck)
	{
		r.POST("/saveResource", api.SaveResource)                //2
		r.POST("/deleteResource", api.DeleteResource)            //0
		r.GET("/getResourceInfo", api.GetResourceInfo)           //0
		r.GET("/getImageList", api.GetImageList)                 //2
		r.GET("/listResource", api.ListResource)                 //0
		r.GET("/changeResourceStatus", api.ChangeResourceStatus) //0
	}
}

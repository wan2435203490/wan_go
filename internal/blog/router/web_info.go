package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, webInfoRouter)
}

func webInfoRouter(v1 *gin.RouterGroup) {
	api := apis.WebInfoApi{}
	r := v1.Group("/webInfo")
	{
		auth := r.Group("").Use(middleware.LoginCheck)
		{
			auth.POST("/updateWebInfo", api.UpdateWebInfo)
			auth.POST("/saveResourcePath", api.SaveResourcePath)
			auth.POST("/saveFriend", api.SaveFriend) //2
			auth.GET("/deleteResourcePath", api.DeleteResourcePath)
			auth.POST("/updateResourcePath", api.UpdateResourcePath)
			auth.POST("/saveFunny", api.SaveFunny)         //2
			auth.POST("/saveLovePhoto", api.SaveLovePhoto) //2
			auth.GET("/deleteTreeHole", api.DeleteTreeHole)
			auth.POST("/saveSort", api.SaveSort)
			auth.GET("/deleteSort", api.DeleteSort)
			auth.POST("/updateSort", api.UpdateSort)
			auth.POST("/saveLabel", api.SaveLabel)
			auth.GET("/deleteLabel", api.DeleteLabel)
			auth.POST("/updateLabel", api.UpdateLabel)
		}

		r.GET("/getWebInfo", api.GetWebInfo)
		r.GET("/getAdmire", api.GetAdmire)
		r.GET("/getSortInfo", api.GetSortInfo)
		r.GET("/getWifeJson", api.GetWifeJson)
		r.GET("/listResourcePath", api.ListResourcePath)
		r.GET("/listFunny", api.ListFunny)
		r.GET("/listCollect", api.ListCollect)
		r.GET("/listAdminLovePhoto", api.ListAdminLovePhoto)
		r.POST("/saveTreeHole", api.SaveTreeHole)
		r.GET("/listTreeHole", api.ListTreeHole)
		r.GET("/listSort", api.ListSort)
		r.GET("/listLabel", api.ListLabel)
		r.GET("/listSortAndLabel", api.ListSortAndLabel)
	}
}

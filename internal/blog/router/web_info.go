package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, webInfoRouter)
}

func webInfoRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.WebInfoApi{}
	r := v1.Group("/webInfo")
	{
		r.GET("", api.GetWebInfo)
		r.GET("/resourcePath/list", api.ListResourcePath)
		r.GET("/sort/list", api.ListSort)
		r.GET("/sort", api.GetSortInfo)
		r.GET("/label/list", api.ListLabel)
		r.POST("/treeHole", api.InsertTreeHole)
		r.GET("/treeHole/list", api.ListTreeHole)
		r.GET("/getAdmire", api.GetAdmire)
		r.GET("/getWaifuJson", api.GetWaifuJson)
		//r.GET("/listFunny", api.ListFunny)
		r.GET("/listCollect", api.ListCollect)
		r.GET("/listAdminLovePhoto", api.ListAdminLovePhoto)
		r.GET("/listSortAndLabel", api.ListSortAndLabel)

		auth1 := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(1))
		{
			auth1.PUT("", api.UpdateWebInfo)
			auth1.POST("/resourcePath", api.InsertResourcePath)
			auth1.DELETE("/resourcePath", api.DeleteResourcePath)
			auth1.PUT("/resourcePath", api.UpdateResourcePath)
			auth1.DELETE("/treeHole", api.DeleteTreeHole)
			auth1.POST("/sort", api.InsertSort)
			auth1.DELETE("/sort", api.DeleteSort)
			auth1.PUT("/sort", api.UpdateSort)
			auth1.POST("/label", api.InsertLabel)
			auth1.DELETE("/label", api.DeleteLabel)
			auth1.PUT("/label", api.UpdateLabel)
		}
		//auth2 := auth.Use(middleware.AuthCheckRole(2))
		//{
		//
		//}
		auth3 := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(3))
		{
			auth3.PUT("/saveFriend", api.SaveFriend)
			auth3.PUT("/saveFunny", api.SaveFunny)
			auth3.PUT("/saveLovePhoto", api.SaveLovePhoto)
		}
	}
}

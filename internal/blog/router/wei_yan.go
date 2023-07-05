package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, weiYanRouter)
}

func weiYanRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.WeiYanApi{}
	r := v1.Group("/weiYan")
	{
		auth := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(3))
		{
			auth.POST("/news", api.InsertNews)
			auth.POST("", api.InsertWeiYan)
			auth.DELETE("", api.DeleteWeiYan)
		}

		r.GET("/listNews", api.ListNews)
		r.GET("/list", api.ListWeiYan)
	}
}

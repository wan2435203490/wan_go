package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, resourceRouter)
}

func resourceRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.ResourceApi{}
	r := v1.Group("/resource")
	{
		auth0 := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(1))
		{
			auth0.DELETE("", api.DeleteResource)
			auth0.GET("", api.GetResourceInfo)
			auth0.GET("/list", api.ListResource)
			auth0.PUT("/changeResourceStatus", api.ChangeResourceStatus)
		}

		auth3 := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(3))
		{
			auth3.POST("", api.InsertResource)
			auth3.GET("/getImageList", api.GetImageList)
		}

	}
}

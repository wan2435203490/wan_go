package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, familyRouter)
}

func familyRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.FamilyApi{}
	r := v1.Group("/family")
	{
		auth1 := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(1))
		{
			auth1.DELETE("", api.DeleteFamily)
			auth1.GET("/list", api.ListFamily)
			auth1.PUT("/changeLoveStatus", api.ChangeLoveStatus)
		}
		auth3 := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(3))
		{
			auth3.POST("", api.InsertFamily)
			auth3.GET("", api.GetFamily)
		}

		r.GET("/getAdminFamily", api.GetAdminFamily)
		r.GET("/listRandomFamily", api.ListRandomFamily)
	}
}

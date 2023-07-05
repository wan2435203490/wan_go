package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, qiniuRouter)
}

func qiniuRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.QiniuApi{}
	r := v1.Group("/qiniu").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(3))
	{
		r.GET("/getUpToken", api.GetUpToken)
		//r.GET("/testUpload", api.TestUpload)
	}
}

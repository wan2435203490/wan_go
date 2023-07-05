package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, commentRouter)
}

func commentRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.CommentApi{}
	r := v1.Group("/comment")
	{
		auth := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(3))
		{
			auth.POST("", api.InsertComment)
			auth.DELETE("", api.DeleteComment)
		}

		r.GET("/count", api.GetCommentCount)
		r.GET("/list", api.ListComment)
	}
}

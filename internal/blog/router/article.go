package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, articleRouter)
}

func articleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.ArticleApi{}
	r := v1.Group("/article")
	{
		auth := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(2))
		{
			auth.POST("", api.InsertArticle)
			auth.DELETE("", api.DeleteArticle)
			auth.PUT("", api.UpdateArticle)
		}

		r.GET("/list", api.ListArticle)
		r.GET("", api.GetArticle)
	}
}

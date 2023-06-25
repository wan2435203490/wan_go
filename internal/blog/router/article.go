package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, articleRouter)
}

func articleRouter(v1 *gin.RouterGroup) {
	api := apis.ArticleApi{}
	r := v1.Group("/article")
	{
		auth := r.Group("").Use(middleware.LoginCheck)
		{
			auth.POST("/saveArticle", api.SaveArticle)     //1
			auth.GET("/deleteArticle", api.DeleteArticle)  //1
			auth.POST("/updateArticle", api.UpdateArticle) //1
		}

		r.GET("/listArticle", api.ListArticle)
		r.GET("/getArticleById", api.GetArticleById)
	}
}

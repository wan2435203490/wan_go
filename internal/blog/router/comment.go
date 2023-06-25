package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
)

func init() {
	routerCheckRole = append(routerCheckRole, commentRouter)
}

func commentRouter(v1 *gin.RouterGroup) {
	api := apis.CommentApi{}
	r := v1.Group("/comment")
	{
		auth := r.Group("") //.Use(middleware.LoginCheck)
		{
			auth.POST("/saveComment", api.SaveComment)
			auth.GET("/deleteComment", api.DeleteComment)
		}

		r.GET("/getCommentCount", api.GetCommentCount)
		r.GET("/listComment", api.ListComment)
	}
}

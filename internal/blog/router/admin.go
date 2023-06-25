package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, adminRouter)
}

func adminRouter(v1 *gin.RouterGroup) {

	api := apis.AdminApi{}
	r := v1.Group("/admin").Use(middleware.LoginCheck)
	{
		r.GET("/user/list", api.ListUser)
		r.GET("/user/changeUserStatus", api.ChangeUserStatus)
		r.GET("/user/changeUserAdmire", api.ChangeUserAdmire)
		r.GET("/user/changeUserType", api.ChangeUserType)

		r.GET("/webInfo/getAdminWebInfo", api.GetAdminWebInfo)

		r.GET("/article/user/list", api.ListUserArticle) //1
		r.GET("/article/boss/list", api.ListBossArticle)
		r.GET("/article/changeArticleStatus", api.ChangeArticleStatus) //1
		r.GET("/article/getArticleById", api.GetArticleByIdForUser)    //1

		r.GET("/comment/user/deleteComment", api.UserDeleteComment) //1
		r.GET("/comment/boss/deleteComment", api.BossDeleteComment)
		r.GET("/comment/user/list", api.ListUserComment) //1
		r.GET("/comment/boss/list", api.ListBossComment)

		r.GET("/treeHole/boss/list", api.ListBossTreeHole)
	}
}

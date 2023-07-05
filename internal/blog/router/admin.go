package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, adminRouter)
}

// todo 抽出一个新的服务
func adminRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	api := apis.AdminApi{}
	r := v1.Group("/admin")
	{
		auth1 := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(1))
		{
			auth1.GET("/user/list", api.ListUser)
			auth1.PUT("/user/changeUserStatus", api.ChangeUserStatus)
			auth1.PUT("/user/changeUserAdmire", api.ChangeUserAdmire)
			auth1.PUT("/user/changeUserType", api.ChangeUserType)
			auth1.GET("/webInfo/getAdminWebInfo", api.GetAdminWebInfo)
			auth1.GET("/treeHole/boss/list", api.ListBossTreeHole)
		}
		auth2 := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(2))
		{
			auth2.GET("/article/list", api.ListArticle)
			auth2.PUT("/article/changeArticleStatus", api.ChangeArticleStatus)
			auth2.GET("/article/getArticleById", api.GetArticleByIdForUser)
			auth2.DELETE("/comment/delete", api.DeleteComment)
			auth2.GET("/comment/list", api.ListComment)
		}

	}
}

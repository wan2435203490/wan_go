package router


import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/landlord/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, chatRouter)
}

func chatRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	api := apis.ChatApi{}
	r := v1.Group("/chat").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckUser())
	{
		r.POST("", api.Chat)
	}

}

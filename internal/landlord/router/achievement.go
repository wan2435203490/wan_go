package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/landlord/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, achievementRouter)
}

func achievementRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	api := apis.AchievementApi{}
	r := v1.Group("/achievement").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckUser())
	{
		r.GET("/:userId", api.GetAchievementByUserId)
	}
}

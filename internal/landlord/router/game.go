package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/landlord/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, gameRouter)
}

func gameRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	api := apis.GameApi{}
	r := v1.Group("/game").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckUser())
	{
		r.POST("/ready", api.Ready)
		r.POST("/unReady", api.UnReady)
		r.POST("/bid", api.Bid)
		r.POST("/play", api.Play)
		r.POST("/pass", api.Pass)
		r.POST("/give", api.Give)
	}

}

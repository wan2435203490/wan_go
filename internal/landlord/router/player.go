package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/landlord/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, playerRouter)
}

func playerRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	api := apis.PlayerApi{}
	r := v1.Group("/player").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckUser())
	{
		r.GET("/cards", api.Cards)
		r.GET("/round", api.Round)
		r.GET("/ready", api.PlayerReady)
		r.GET("/pass", api.PlayerPass)
		r.GET("/bidding", api.Bidding)
	}

}

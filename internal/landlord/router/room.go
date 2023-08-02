package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/landlord/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, roomRouter)
}

func roomRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	api := apis.RoomApi{}
	r := v1.Group("/room").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckUser())
	{
		r.GET("", api.Rooms)
		r.GET("/:id", api.GetById)
		r.POST("", api.Create)
		r.POST("/join", api.Join)
		r.POST("/exit", api.Exit)
	}

}

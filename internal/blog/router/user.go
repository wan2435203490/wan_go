package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, userRouter)
}

func userRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.UserApi{}
	r := v1.Group("/user")
	{
		auth := r.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole(3))
		{
			auth.PUT("", api.UpdateUserInfo)
			auth.PUT("/updateSecretInfo", api.UpdateSecretInfo)
			auth.GET("/getUserByUsername", api.GetUserByUsername)
			auth.GET("/getInfo", api.GetInfo)
			auth.POST("/refreshToken", authMiddleware.RefreshHandler)
		}

		r.GET("/getCaptchaForBind", api.GetCaptchaForBind)
		r.GET("/getCaptcha", api.GetCaptcha)
		r.GET("/getCaptchaImg", api.GetCaptchaImg)
		r.GET("/logout", api.Logout)
		r.POST("/register", api.Register)
		r.POST("/login", authMiddleware.LoginHandler)
		r.POST("/token", api.LoginByToken)
		r.GET("/getCaptchaForForgetPassword", api.GetCaptchaForForgetPassword)
		r.PUT("/updateForForgetPassword", api.UpdateForForgetPassword)
	}
}

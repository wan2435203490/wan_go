package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, userRouter)
}

func userRouter(v1 *gin.RouterGroup) {
	api := apis.UserApi{}
	r := v1.Group("/user")
	{
		auth := r.Group("").Use(middleware.LoginCheck)
		{
			auth.GET("/logout", api.Logout)                       //2
			auth.POST("/updateUserInfo", api.UpdateUserInfo)      //2
			auth.GET("/getCode", api.GetCaptcha)                  //2
			auth.GET("/getCodeForBind", api.GetCaptchaForBind)    //2
			auth.POST("/updateSecretInfo", api.UpdateSecretInfo)  //2
			auth.GET("/getUserByUsername", api.GetUserByUsername) //2
		}
		r.POST("/register", api.Register)
		r.POST("/login", api.Login)
		r.POST("/token", api.LoginByToken)
		r.GET("/getCodeForForgetPassword", api.GetCaptchaForForgetPassword)
		r.POST("/updateForForgetPassword", api.UpdateForForgetPassword)
	}
}

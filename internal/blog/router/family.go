package router

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/apis"
	"wan_go/pkg/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, familyRouter)
}

func familyRouter(v1 *gin.RouterGroup) {
	api := apis.FamilyApi{}
	r := v1.Group("/family")
	{
		auth := r.Group("").Use(middleware.LoginCheck)
		{
			auth.POST("/saveFamily", api.SaveFamily)            //2
			auth.GET("/deleteFamily", api.DeleteFamily)         //0
			auth.GET("/getFamily", api.GetFamily)               //2
			auth.GET("/listFamily", api.ListFamily)             //0
			auth.GET("/changeLoveStatus", api.ChangeLoveStatus) //0
		}

		r.GET("/getAdminFamily", api.GetAdminFamily)
		r.GET("/listRandomFamily", api.ListRandomFamily)
	}
}

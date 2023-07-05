package middleware

import (
	"github.com/gin-gonic/gin"
	"wan_go/pkg/common/constant"
	"wan_go/sdk/pkg/jwtauth"
	"wan_go/sdk/pkg/response"
)

// AuthCheckRole 权限检查中间件
// roleId 1 站长
// roleId 2 管理员
// roleId 3 普通角色
// roleId 当前要求的最高角色Id 比如传的2，只有v[constant.ClaimsRoleId]为1，2可以进当前router
func AuthCheckRole(roleId int) gin.HandlerFunc {
	//暂不实现casbin 没必要
	return func(c *gin.Context) {
		//log := api.GetRequestLogger(c)
		data, _ := c.Get(constant.JWTPayloadKey)
		v := data.(jwtauth.MapClaims)
		//e := sdk.Runtime.GetCasbinKey(c.Request.Host)
		//var res, casbinExclude bool
		//var err error
		//检查权限
		if int(v[constant.ClaimsRoleId].(float64)) > roleId {
			//res = true
			response.Error(c, 500, nil, "权限不足")
			c.Abort()
			return
		}

		c.Next()

		//for _, i := range CasbinExclude {
		//	if util.KeyMatch2(c.Request.URL.Path, i.Url) && c.Request.Method == i.Method {
		//		casbinExclude = true
		//		break
		//	}
		//}
		//if casbinExclude {
		//	log.Infof("Casbin exclusion, no validation method:%s path:%s", c.Request.Method, c.Request.URL.Path)
		//	c.Next()
		//	return
		//}
		//res, err = e.Enforce(v[constant.ClaimsRoleId], c.Request.URL.Path, c.Request.Method)
		//if err != nil {
		//	log.Errorf("AuthCheckRole error:%s method:%s path:%s", err, c.Request.Method, c.Request.URL.Path)
		//	response.Error(c, 500, err, "")
		//	return
		//}
		//
		//if res {
		//	log.Infof("isTrue: %v role: %s method: %s path: %s", res, v[constant.ClaimsRoleId], c.Request.Method, c.Request.URL.Path)
		//	c.Next()
		//} else {
		//	log.Warnf("isTrue: %v role: %s method: %s path: %s message: %s", res, v[constant.ClaimsRoleId], c.Request.Method, c.Request.URL.Path, "当前request无权限，请管理员确认！")
		//	c.JSON(http.StatusOK, gin.H{
		//		"code": 403,
		//		"msg":  "对不起，您没有该接口访问权限，请联系管理员",
		//	})
		//	c.Abort()
		//	return
		//}

	}
}

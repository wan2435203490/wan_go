package middleware

import (
	"time"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/middleware/handler"
	jwt "wan_go/sdk/pkg/jwtauth"
)

// AuthInit jwt验证new
func AuthInit() (*jwt.GinJWTMiddleware, error) {
	timeout := time.Hour
	if config.Config.Application.Mode == "dev" {
		timeout *= time.Duration(876010)
	} else if config.Config.Jwt.Expire != 0 {
		timeout *= time.Duration(config.Config.Jwt.Expire) * 24
	}

	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "wan zone",
		Key:             []byte(config.Config.Jwt.Secret),
		Timeout:         timeout,
		MaxRefresh:      time.Hour,
		PayloadFunc:     handler.PayloadFunc,     //第一次登录获取token
		IdentityHandler: handler.IdentityHandler, //后面每次进来解析jwt claimsMap 返回成map
		Authenticator:   handler.Authenticator,   //登录验证user role captcha
		Authorizator:    handler.Authorizator,    //将jwt的内容set到gin.context
		Unauthorized:    handler.Unauthorized,    //验证失败调用
		TokenLookup:     "header: Authorization", //token 解析
		TokenHeadName:   "WanBlog",               //header Authorization prefix
		TimeFunc:        time.Now,                //token expire
	})

}

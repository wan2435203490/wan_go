package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/sdk/pkg"
	"wan_go/sdk/pkg/captcha"
	jwt "wan_go/sdk/pkg/jwtauth"
	"wan_go/sdk/pkg/response"
)

// PayloadFunc LoginHandler
func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v[constant.User].(blog.User)
		r, _ := v[constant.Role].(blog.Role)
		return jwt.MapClaims{
			jwt.Identity:    u.ID,
			jwt.UserName:    u.UserName,
			jwt.Avatar:      u.Avatar,
			jwt.Password:    u.Password,
			jwt.PhoneNumber: u.PhoneNumber,
			jwt.Email:       u.Email,
			jwt.RoleId:      r.ID,
			jwt.RoleName:    r.Name,
			jwt.DataScope:   r.DataScope,
		}
	}
	return jwt.MapClaims{}
}

// IdentityHandler Get From PayloadFunc
func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return map[string]interface{}{
		jwt.Identity:    claims[jwt.Identity],
		jwt.UserName:    claims[jwt.UserName],
		jwt.Avatar:      claims[jwt.Avatar],
		jwt.Password:    claims[jwt.Password],
		jwt.PhoneNumber: claims[jwt.PhoneNumber],
		jwt.Email:       claims[jwt.Email],
		jwt.RoleId:      claims[jwt.RoleId],
		jwt.RoleName:    claims[jwt.RoleName],
		jwt.DataScope:   claims[jwt.DataScope],
	}
}

// Authenticator 获取token
//
//	@Summary		登陆
//	@Description	获取token
//	@Description	LoginHandler can be used by clients to get a jwt token.
//	@Description	Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
//	@Description	Reply will be of the form {"token": "TOKEN"}.
//	@Description	dev mode：It should be noted that all fields cannot be empty, and a value of 0 can be passed in addition to the account password
//	@Description	注意：开发模式：需要注意全部字段不能为空，账号密码外可以传入0值
//	@Tags			登陆
//	@Accept			application/json
//	@Product		application/json
//	@Param			account	body		Login	true	"account"
//	@Success		200		{string}	string	"{"code": 200, "expire": "2019-08-07T12:45:48+08:00", "token": ".eyJleHAiOjE1NjUxNTMxNDgsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU2NTE0OTU0OH0.-zvzHvbg0A" }"
//	@Router			/api/v1/login [post]
func Authenticator(c *gin.Context) (interface{}, error) {
	log := api.GetRequestLogger(c)
	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db error, %s", err.Error())
		response.Error(c, 500, err, "数据库连接获取失败")
		return nil, jwt.ErrFailedAuthentication
	}

	var loginVals Login
	var status = "2"
	var msg = "登录成功"
	var username = ""
	//defer func() {
	//	LoginLogToDB(c, status, msg, username)
	//}()

	if err = c.Bind(&loginVals); err != nil {
		username = loginVals.UserName
		msg = "数据解析失败"
		status = "1"

		return nil, jwt.ErrMissingLoginValues
	}
	if config.Config.Application.Mode != "dev" && !loginVals.IsAdmin {
		if !captcha.Verify(loginVals.UUID, loginVals.Captcha, true) {
			username = loginVals.UserName
			msg = "验证码错误"
			status = "1"

			return nil, jwt.ErrInvalidVerificationode
		}
	}
	user, role, e := loginVals.GetUser(db)
	if e == nil {
		username = loginVals.UserName

		return map[string]interface{}{constant.User: user, constant.Role: role}, nil
	} else {
		msg = "登录失败"
		status = "1"
		log.Warnf("%s login failed!", loginVals.UserName)
	}
	log.Info(status, username, msg)
	return nil, jwt.ErrFailedAuthentication
}

// LogOut
//
//	@Summary		退出登录
//	@Description	获取token
//
// LoginHandler can be used by clients to get a jwt token.
// Reply will be of the form {"token": "TOKEN"}.
//
//	@Accept			application/json
//	@Product		application/json
//	@Success		200	{string}	string	"{"code": 200, "msg": "成功退出系统" }"
//	@Router			/logout [post]
//	@Security		Bearer
func LogOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出成功",
	})
}

func Authorizator(data interface{}, c *gin.Context) bool {

	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v[constant.User].(blog.User)
		r, _ := v[constant.Role].(blog.Role)
		c.Set(constant.ClaimsIdentity, u.ID)
		c.Set(constant.ClaimsUserName, u.UserName)
		c.Set(constant.ClaimsPassword, u.Password)
		c.Set(constant.ClaimsEmail, u.Email)
		c.Set(constant.ClaimsPhoneNumber, u.PhoneNumber)
		c.Set(constant.ClaimsRoleId, r.ID)
		c.Set(constant.ClaimsRoleName, r.Name)
		c.Set(constant.ClaimsDataScope, r.DataScope)
		return true
	}
	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  message,
	//})
	//res := &r.Response{}
	//res.Message = message
	//res.Status = r.ErrorStatus
	//res.Captcha = code
	//
	//c.JSON(code, res)

	response.Error(c, code, nil, message)
}

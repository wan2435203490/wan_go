package user

import (
	"fmt"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"

	"github.com/gin-gonic/gin"

	"wan_go/sdk/pkg"
	jwt "wan_go/sdk/pkg/jwtauth"
)

func ExtractClaims(c *gin.Context) jwt.MapClaims {
	claims, exists := c.Get(constant.JWTPayloadKey)
	if !exists || claims == nil {
		return make(jwt.MapClaims)
	}

	return claims.(jwt.MapClaims)
}

func RemoveToken(c *gin.Context) {
	c.Set(constant.JWTPayloadKey, nil)
}

// Get constant.ClaimsIdentity ...
func Get(c *gin.Context, key string) interface{} {
	data := ExtractClaims(c)
	if data[key] != nil {
		return data[key]
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " Get 缺少 " + key)
	return nil
}

func GetUserInfo(c *gin.Context) (userId int32, userName string) {
	data := ExtractClaims(c)
	if data[constant.ClaimsIdentity] != nil {
		userId = int32((data[constant.ClaimsIdentity]).(float64))
		userName = (data[constant.ClaimsUserName]).(string)
		return
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserId 缺少 identity")
	return 0, ""
}

func GetUser(c *gin.Context, user *blog.User) {
	data := ExtractClaims(c)
	if data[constant.ClaimsIdentity] != nil {
		user.ID = int32((data[constant.ClaimsIdentity]).(float64))
		user.UserName = (data[constant.ClaimsUserName]).(string)
		user.Avatar = (data[constant.ClaimsAvatar]).(string)
		user.Email = (data[constant.ClaimsEmail]).(string)
		user.PhoneNumber = (data[constant.ClaimsPhoneNumber]).(string)
		return
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUser 缺少 identity")
}

func GetUserId32(c *gin.Context) int32 {
	data := ExtractClaims(c)
	if data[constant.ClaimsIdentity] != nil {
		return int32((data[constant.ClaimsIdentity]).(float64))
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserId32 缺少 identity")
	return 0
}

func GetUserId(c *gin.Context) int {
	data := ExtractClaims(c)
	if data[constant.ClaimsIdentity] != nil {
		return int((data[constant.ClaimsIdentity]).(float64))
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserId 缺少 identity")
	return 0
}

func GetUserIdStr(c *gin.Context) string {
	data := ExtractClaims(c)
	if data[constant.ClaimsIdentity] != nil {
		return pkg.Int64ToString(int64((data[constant.ClaimsIdentity]).(float64)))
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserIdStr 缺少 identity")
	return ""
}

func GetUserName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data[constant.ClaimsUserName] != nil {
		return (data[constant.ClaimsUserName]).(string)
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserName 缺少 username")
	return ""
}

func GetUserPassword(c *gin.Context) string {
	data := ExtractClaims(c)
	if data[constant.ClaimsPassword] != nil {
		return (data[constant.ClaimsPassword]).(string)
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserPassword 缺少 password")
	return ""
}

func GetEmail(c *gin.Context) string {
	data := ExtractClaims(c)
	if data[constant.ClaimsEmail] != nil {
		return (data[constant.ClaimsEmail]).(string)
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetEmail 缺少 email")
	return ""
}

func GetUserPhoneNumber(c *gin.Context) string {
	data := ExtractClaims(c)
	if data[constant.ClaimsPhoneNumber] != nil {
		return (data[constant.ClaimsPhoneNumber]).(string)
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserPhoneNumber 缺少 phoneNumber")
	return ""
}

func GetRoleName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data[constant.ClaimsRoleName] != nil {
		return (data[constant.ClaimsRoleName]).(string)
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetRoleName 缺少 roleName")
	return ""
}

// IsAdmin 站长或管理员 返回true
func IsAdmin(c *gin.Context) bool {
	roleId := int8(GetRoleId(c))
	return roleId < blog_const.UserRoleOperator.Code
}

func GetRoleId(c *gin.Context) int {
	data := ExtractClaims(c)
	if data[constant.ClaimsRoleId] != nil {
		i := int((data[constant.ClaimsRoleId]).(float64))
		return i
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetRoleId 缺少 roleId")
	return 0
}

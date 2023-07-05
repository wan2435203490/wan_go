package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	r "wan_go/pkg/common/response"
	"wan_go/sdk/pkg/response"
)

//var mu sync.Mutex

//func Test1(c *gin.Context) {
//	//abort后不会执行后面的中间件
//	c.Abort()
//	fmt.Println("test1")
//
//}
//
//func Test2(c *gin.Context) {
//	c.Abort()
//	fmt.Println("test2")
//}

//func LoginCheck(c *gin.Context) {
//
//	tokens := c.Request.Header[blog_const.TOKEN_HEADER]
//	if abort(c, len(tokens) == 0 || len(tokens[0]) == 0, &r.NOT_LOGIN) {
//		return
//	}
//
//	token := tokens[0]
//	cache.SetToken(token)
//	user := cache.GetUser()
//	//user := &blog.User{ID: 1, UserStatus: true, UserType: 0, UserName: "wan"}
//	//cache.SetUser(user)
//	if abort(c, user == nil, &r.LOGIN_EXPIRED) {
//		return
//	}
//
//	userId := utils.Int32ToString(user.ID)
//	isUser, isAdmin, key := userAdmin(token, userId)
//
//	var loginCheck int8 = 2
//	if isUser {
//		cond := loginCheck == blog_const.UserRoleAdmin.Captcha || loginCheck == blog_const.UserRoleDev.Captcha
//		if abort(c, cond, &r.FAIL_ADMIN) {
//			return
//		}
//	} else if isAdmin {
//		cond := loginCheck == blog_const.UserRoleAdmin.Captcha && int(user.ID) != blog_const.ADMIN_USER_ID
//		if abort(c, cond, &r.FAIL_ADMIN) {
//			return
//		}
//	} else {
//		c.AbortWithStatusJSON(r.NOT_LOGIN.CodeMsg())
//		return
//	}
//
//	if abort(c, loginCheck < user.UserType, &r.FAIL_PERSSION) {
//		return
//	}
//
//	get, b := cache.Get(key)
//	flag1 := !b || get.(string) == ""
//	if flag1 {
//		mu.Lock()
//		defer mu.Unlock()
//		get, b = cache.Get(key)
//		flag2 := !b || get.(string) == ""
//
//		if flag2 {
//			cache.SetExpire(token, user, blog_const.TOKEN_EXPIRE)
//
//			cache.SetExpire(key, token, blog_const.TOKEN_EXPIRE)
//			if isUser {
//				cache.SetExpire(blog_const.USER_TOKEN_INTERVAL+userId, token, blog_const.TOKEN_INTERVAL)
//			} else if isAdmin {
//				cache.SetExpire(blog_const.ADMIN_TOKEN_INTERVAL+userId, token, blog_const.TOKEN_INTERVAL)
//			}
//		}
//	}
//}

func RemoveLocalToken(c *gin.Context) {
	cache.RemoveToken()
}

func userAdmin(token, userId string) (isUser bool, isAdmin bool, key string) {
	isUser = strings.Contains(token, blog_const.USER_ACCESS_TOKEN)
	isAdmin = strings.Contains(token, blog_const.ADMIN_ACCESS_TOKEN)

	if isUser {
		key = blog_const.USER_TOKEN_INTERVAL + userId
	} else if isAdmin {
		key = blog_const.ADMIN_TOKEN_INTERVAL + userId
	}

	return
}

func abort(c *gin.Context, cond bool, cm *r.CodeMsg) bool {
	if cond {
		//res := &r.Response{}
		//res.Status = r.ErrorStatus
		//res.Captcha, res.Message = cm.CodeMsg()
		//c.AbortWithStatusJSON(res.Captcha, res)
		response.Error(c, cm.Code, nil, cm.Msg)
		return true
	}
	return false
}

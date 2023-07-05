package cache

//单机cache
import (
	"github.com/patrickmn/go-cache"
	"github.com/timandy/routine"
	"time"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/utils"
)

var (
	//默认过期时间为 6 小时的缓存，每 12 小时清除一次过期key
	cc    = cache.New(6*time.Hour, 12*time.Hour)
	local = routine.NewThreadLocal()
)

func SetToken(token string) {
	local.Set(token)
}

func RemoveToken() {
	local.Remove()
}

func Token() string {
	get := local.Get()
	if get == nil {
		return ""
	}
	return get.(string)
}

// Set default cache.DefaultExpiration
func Set(key string, value any) {
	cc.Set(key, value, cache.DefaultExpiration)
}

func SetExpire(key string, value any, d time.Duration) {
	//第三个参数为该key过期时间，大于0时生效 default时取的new函数的第一个值
	cc.Set(key, value, cache.DefaultExpiration)
}

func SetCaptchaExpire(userId string, flag int, captcha string) {
	SetExpire(utils.CaptchaKey(userId, flag), captcha, time.Minute*5)
}

func SetCaptchaBindExpire(userId, place string, flag int, captcha string) {
	SetExpire(utils.CaptchaBindKey(userId, place, flag), captcha, time.Minute*5)
}

func Get(key string) (any, bool) {
	return cc.Get(key)
}

func GetString(key string) (string, bool) {
	get, b := cc.Get(key)
	if get == nil {
		return "", false
	}
	return get.(string), b
}

func Delete(key string) {
	cc.Delete(key)
}

func GetAdminUser() *blog.User {
	if get, b := Get(blog_const.ADMIN); b {
		return get.(*blog.User)
	}
	return nil
}

func GetAdminUserId() int {
	if admin := GetAdminUser(); admin != nil {
		return int(admin.ID)
	}
	return -1
}

//
//func SetUser(user *blog.User) {
//	Set(Token(), user)
//}
//
//func GetUser() *blog.User {
//	if get, b := Get(Token()); b {
//		return get.(*blog.User)
//	}
//	return nil
//}
//
//func GetUserId() int {
//	if user := GetUser(); user != nil {
//		return int(user.ID)
//	}
//	return -1
//}
//
//func GetUserIdStr() string {
//	if user := GetUser(); user != nil {
//		return utils.Int32ToString(user.ID)
//	}
//	return ""
//}
//
//func GetUserName() string {
//	user := GetUser()
//	if user != nil {
//		return user.UserName
//	}
//	return ""
//}
//
//func GetWebInfo() *blog.WebInfo {
//	if get, b := Get(blog_const.WEB_INFO); b {
//		return get.(*blog.WebInfo)
//	}
//	var err error
//	ws := service.NewWebInfo(nil)
//	var webInfos []blog.WebInfo
//	if err = ws.List(&webInfos); err != nil {
//		panic(err)
//	}
//	return &webInfos[0]
//}

func GetWebName() string {
	//webInfo := GetWebInfo()
	//if webInfo == nil {
	//	return config.Config.Blog.Name
	//}
	//return webInfo.WebName
	return config.Config.Blog.Name
}

//func CanSendEmail(user *blog.User) bool {
//	return user != nil && int(user.ID) != GetUserId() && len(user.Email) > 0
//}

package cache

//单机cache
import (
	"github.com/patrickmn/go-cache"
	"time"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	//blogVO "wan_go/pkg/vo/blog"
)

var (
	//默认过期时间为 30 分钟的缓存，每 1 小时清除一次过期key
	cc = cache.New(30*time.Minute, 1*time.Hour)
)

// Set default cache.DefaultExpiration
func Set(key string, value any) {
	cc.Set(key, value, cache.DefaultExpiration)
}

func SetExpire(key string, value any, d time.Duration) {
	//第三个参数为该key过期时间，大于0时生效 default时取的new函数的第一个值
	cc.Set(key, value, d)
}

func Get(key string) (any, bool) {
	return cc.Get(key)
}

func GetString(key string) (string, bool) {
	get, b := cc.Get(key)
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

func GetWebInfo() *blog.WebInfo {
	if get, b := Get(blog_const.WEB_INFO); b {
		return get.(*blog.WebInfo)
	}
	return nil
}

func GetWebName() string {
	webInfo := GetWebInfo()
	if webInfo == nil {
		return config.Config.Blog.Name
	}
	return webInfo.WebName
}
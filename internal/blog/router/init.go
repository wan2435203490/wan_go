package router

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog/db_common"
	"wan_go/pkg/common/db/mysql/blog/db_family"
	"wan_go/pkg/common/db/mysql/blog/db_user"
	"wan_go/pkg/common/db/mysql/blog/db_web_info"
	"wan_go/pkg/common/middleware"
	"wan_go/pkg/common/prometheus"
)

var (
	routerCheckRole = make([]func(v1 *gin.RouterGroup), 0)
	store           = cookie.NewStore([]byte(config.Config.Session.Secret))
	//暂时不用redis
	//store, _ = redis.NewStore(10, "tcp", "localhost:6379", "", []byte(config.Config.Session.Secret))
)

func InitMiddleware(engine *gin.Engine) {
	engine.Use(sessions.Sessions(config.Config.Session.Name, store))
	//engine.Use(middleware.WithCors)
	engine.Use(middleware.NoCache, middleware.Options, middleware.Secure)
	//engine.Use(RequestId(TrafficKey), apis.SetRequestLogger)
	//engine.Use(CustomError, NoCache)
	//todo logs permission
}

func Init(engine *gin.Engine) {
	initCache()
	initRouter(engine)

	engine.Use(middleware.RemoveLocalToken)
}

func initCache() {
	webInfos, err := db_web_info.List()
	if err != nil || len(webInfos) == 0 {
		log.Println("web info is empty")
	} else {
		cache.Set(blog_const.WEB_INFO, webInfos[0])
	}

	sortInfos := db_common.GetSortInfo()
	if len(*sortInfos) == 0 {
		log.Println("sortInfos is empty")
	} else {
		cache.Set(blog_const.SORT_INFO, sortInfos)
	}

	admin := db_user.GetByUserType(blog_const.USER_TYPE_ADMIN.Code)
	if admin != nil && admin.ID > 0 {
		cache.Set(blog_const.ADMIN, admin)
	} else {
		log.Println("admin is not exist")
	}

	family := db_family.GetByUserId(int(admin.ID))
	cache.Set(blog_const.ADMIN_FAMILY, family)
}

func initRouter(engine *gin.Engine) {

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if config.Config.Prometheus.Enable {
		engine.GET("/metrics", prometheus.PrometheusHandler())
	}
	InitMiddleware(engine)

	checkRoleRouter(engine)
}

// 需要认证的路由示例
func checkRoleRouter(r *gin.Engine) {
	v1 := r.Group("")
	for _, f := range routerCheckRole {
		f(v1)
	}
}

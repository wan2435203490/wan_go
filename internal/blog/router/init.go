package router

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	//_ "github/mwqnice/swag/docs"
	"log"
	"sync"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/db"
	rocksCache "wan_go/pkg/common/db/rocks_cache"
	"wan_go/pkg/common/middleware"
	"wan_go/pkg/common/prometheus"
	"wan_go/sdk"
	jwt "wan_go/sdk/pkg/jwtauth"
)

var (
	routerNoCheckRole = make([]func(v1 *gin.RouterGroup), 0)
	routerCheckRole   = make([]func(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware), 0)
	store             = cookie.NewStore([]byte(config.Config.Session.Secret))
	//暂时不用redis
	//store, _ = redis.NewStore(10, "tcp", "localhost:6379", "", []byte(config.Config.Session.Secret))
)

func InitMiddleware(engine *gin.Engine) {
	// 数据库链接
	engine.Use(middleware.WithContextDb)
	engine.Use(sessions.Sessions(config.Config.Session.Name, store))
	//engine.Use(middleware.WithCors)
	engine.Use(middleware.NoCache, middleware.Options, middleware.Secure)
	//engine.Use(middleware.RequestId(TrafficKey), apis.SetRequestLogger)
	//engine.Use(CustomError, NoCache)
}

func Init(engine *gin.Engine) {
	initDB()
	initCache()
	initRouter(engine)

	engine.Use(middleware.RemoveLocalToken)
}

func initDB() {
	sdk.Runtime.SetDb("*", db.Mysql())
}

func initCache() {

	rocksCache.DelBlogKeys()

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		webInfos, err := rocksCache.GetWebInfo()
		if err != nil || webInfos == nil {
			panic("webInfos is empty")
		}
		wg.Done()
	}()

	go func() {
		sortInfos, err := rocksCache.GetSortInfo()
		if err != nil || sortInfos == nil {
			panic("sortInfos is empty")
		}
		wg.Done()
	}()

	go func() {
		admire, err := rocksCache.GetAdmire()
		if err != nil || admire == nil {
			panic("admire is empty")
		}
		wg.Done()
	}()

	go func() {
		admin, err := rocksCache.GetAdminUser()
		if err != nil || admin == nil {
			panic("admin is no exist")
		}

		family, err := rocksCache.GetAdminFamily(admin.ID)
		if err != nil || family == nil {
			panic("family is no exist")
		}

		wg.Done()
	}()

	wg.Wait()
	fmt.Println("\nrocks cache inited")
}

func initRouter(engine *gin.Engine) {

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if config.Config.Prometheus.Enable {
		engine.GET("/metrics", prometheus.PrometheusHandler())
	}

	// the jwt middleware
	authMiddleware, err := middleware.AuthInit()
	if err != nil {
		log.Fatalf("JWT Init Error, %s", err.Error())
	}

	InitMiddleware(engine)

	checkRoleRouter(engine, authMiddleware)
}

// 需要认证的路由示例
func checkRoleRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := r.Group("/api")
	for _, f := range routerNoCheckRole {
		f(v1)
	}
	v2 := r.Group("/api")
	for _, f := range routerCheckRole {
		f(v2, authMiddleware)
	}
}

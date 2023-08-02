package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"wan_go/pkg/common/api"
	rocksCache "wan_go/pkg/common/db/rocks_cache"
	"wan_go/pkg/utils"

	//_ "github/mwqnice/swag/docs"
	"log"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/middleware"
	"wan_go/pkg/common/prometheus"
	"wan_go/sdk"
	jwt "wan_go/sdk/pkg/jwtauth"
)

var (
	routerNoCheckRole = make([]func(v1 *gin.RouterGroup), 0)
	routerCheckRole   = make([]func(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware), 0)
	//store             = cookie.NewStore([]byte(config.Config.Session.Secret))
)

func InitMiddleware(engine *gin.Engine) {
	// db
	engine.Use(middleware.WithContextDb)
	//engine.Use(sessions.Sessions(config.Config.Session.Name, store))
	engine.Use(middleware.NoCache, middleware.Options, middleware.Secure)
	engine.Use(middleware.RequestId(utils.TrafficKey), api.SetRequestLogger)
	//engine.Use(CustomError)
}

func Init(engine *gin.Engine) {
	initDB()
	initCache()
	initRouter(engine)

	//engine.Use(middleware.RemoveLocalToken)
}

func initDB() {
	sdk.Runtime.SetDb("*", db.Mysql())
}

func initCache() {

	rocksCache.DelLandlordKeys()
	//
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//
	//wg.Wait()
	//fmt.Println("\nrocks-cache inited")
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
	v1 := r.Group("/landlord/api")
	for _, f := range routerNoCheckRole {
		f(v1)
	}
	v2 := r.Group("/landlord/api")
	for _, f := range routerCheckRole {
		f(v2, authMiddleware)
	}
}

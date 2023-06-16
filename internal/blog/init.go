package blog

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/middleware"
	"wan_go/pkg/common/prometheus"
)

var (
	store = cookie.NewStore([]byte(config.Config.Session.Secret))
	//暂时不用redis
	//store, _ = redis.NewStore(10, "tcp", "localhost:6379", "", []byte(config.Config.Session.Secret))
)

func InitMiddleware(engine *gin.Engine) {
	engine.Use(sessions.Sessions(config.Config.Session.Name, store))
	engine.Use(middleware.WithCors, middleware.WithSession, middleware.WithLimit)
	//engine.Use(RequestId(TrafficKey), api.SetRequestLogger)
	//engine.Use(CustomError, NoCache)
	//todo log permission
}

func Init(engine *gin.Engine) {

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if config.Config.Prometheus.Enable {
		engine.GET("/metrics", prometheus.PrometheusHandler())
	}
	InitMiddleware(engine)

	webInfo := engine.Group("/webInfo")
	{
		webInfo.Use(WebInfo)
		webInfo.POST("/updateWebInfo", UpdateWebInfo)
		webInfo.GET("/getWebInfo", UpdateWebInfo)
		webInfo.GET("/getAdmire", UpdateWebInfo)
		webInfo.GET("/getSortInfo", UpdateWebInfo)
		webInfo.GET("/getWaifuJson", UpdateWebInfo)
		webInfo.POST("/saveResourcePath", UpdateWebInfo)
		webInfo.POST("/saveFriend", UpdateWebInfo)
		webInfo.GET("/deleteResourcePath", UpdateWebInfo)
		webInfo.POST("/updateResourcePath", UpdateWebInfo)
		webInfo.POST("/listResourcePath", UpdateWebInfo)
		webInfo.GET("/listFunny", UpdateWebInfo)
		webInfo.GET("/listCollect", UpdateWebInfo)
		webInfo.POST("/saveFunny", UpdateWebInfo)
		webInfo.GET("/listAdminLovePhoto", UpdateWebInfo)
		webInfo.POST("/saveLovePhoto", UpdateWebInfo)
		webInfo.POST("/saveTreeHole", UpdateWebInfo)
		webInfo.GET("/deleteTreeHole", UpdateWebInfo)
		webInfo.GET("/listTreeHole", UpdateWebInfo)
		webInfo.POST("/saveSort", UpdateWebInfo)
		webInfo.GET("/deleteSort", UpdateWebInfo)
		webInfo.POST("/updateSort", UpdateWebInfo)
		webInfo.GET("/listSort", UpdateWebInfo)
		webInfo.POST("/saveLabel", UpdateWebInfo)
		webInfo.GET("/deleteLabel", UpdateWebInfo)
		webInfo.POST("/updateLabel", UpdateWebInfo)
		webInfo.GET("/listLabel", UpdateWebInfo)
		webInfo.GET("/listSortAndLabel", UpdateWebInfo)
	}
}

func Run(prometheusPort int) {
	//go ws.run()

	go func() {
		err := prometheus.StartPromeSrv(prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
}

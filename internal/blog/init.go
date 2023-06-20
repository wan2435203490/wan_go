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
	engine.Use(middleware.WithCors, middleware.WithLimit)
	//engine.Use(RequestId(TrafficKey), api.SetRequestLogger)
	//engine.Use(CustomError, NoCache)
	//todo logs permission
}

func Init(engine *gin.Engine) {

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if config.Config.Prometheus.Enable {
		engine.GET("/metrics", prometheus.PrometheusHandler())
	}
	InitMiddleware(engine)

	engine.Use(Api)

	webInfo := engine.Group("/webInfo")
	{
		webInfo.POST("/updateWebInfo", UpdateWebInfo)
		webInfo.GET("/getWebInfo", GetWebInfo)
		webInfo.GET("/getAdmire", GetAdmire)
		webInfo.GET("/getSortInfo", GetSortInfo)
		webInfo.GET("/getWaifuJson", GetWifeJson)
		webInfo.POST("/saveResourcePath", SaveResourcePath)
		webInfo.POST("/saveFriend", SaveFriend)
		webInfo.GET("/deleteResourcePath", DeleteResourcePath)
		webInfo.POST("/updateResourcePath", UpdateResourcePath)
		webInfo.POST("/listResourcePath", ListResourcePath)
		webInfo.GET("/listFunny", ListFunny)
		webInfo.GET("/listCollect", ListCollect)
		webInfo.POST("/saveFunny", SaveFunny)
		webInfo.GET("/listAdminLovePhoto", ListAdminLovePhoto)
		webInfo.POST("/saveLovePhoto", SaveLovePhoto)
		webInfo.POST("/saveTreeHole", SaveTreeHole)
		webInfo.GET("/deleteTreeHole", DeleteTreeHole)
		webInfo.GET("/listTreeHole", ListTreeHole)
		webInfo.POST("/saveSort", SaveSort)
		webInfo.GET("/deleteSort", DeleteSort)
		webInfo.POST("/updateSort", UpdateSort)
		webInfo.GET("/listSort", ListSort)
		webInfo.POST("/saveLabel", SaveLabel)
		webInfo.GET("/deleteLabel", DeleteLabel)
		webInfo.POST("/updateLabel", UpdateLabel)
		webInfo.GET("/listLabel", ListLabel)
		webInfo.GET("/listSortAndLabel", ListSortAndLabel)
	}
	admin := engine.Group("/admin")
	{
		admin.POST("/user/list", ListUser)
		admin.GET("/user/changeUserStatus", ChangeUserStatus)
		admin.GET("/user/changeUserAdmire", ChangeUserAdmire)
		admin.GET("/user/changeUserType", ChangeUserType)
		admin.GET("/webInfo/getAdminWebInfo", GetAdminWebInfo)
		admin.POST("/article/user/list", ListUserArticle)
		admin.POST("/article/boss/list", ListBossArticle)
		admin.GET("/article/changeArticleStatus", ChangeArticleStatus)
		admin.GET("/article/getArticleById", GetArticleByIdForUser)
		admin.GET("/comment/user/deleteComment", UserDeleteComment)
		admin.GET("/comment/boss/deleteComment", BossDeleteComment)
		admin.POST("/comment/user/list", ListUserComment)
		admin.POST("/comment/boss/list", ListBossComment)
		admin.POST("/treeHole/boss/list", ListBossTreeHole)
	}
	article := engine.Group("/article")
	{
		article.POST("/saveArticle", SaveArticle)
		article.GET("/deleteArticle", DeleteArticle)
		article.POST("/updateArticle", UpdateArticle)
		article.POST("/listArticle", ListArticle)
		article.GET("/getArticleById", GetArticleById)
	}
	comment := engine.Group("/comment")
	{
		comment.POST("/saveComment", SaveComment)
		comment.GET("/deleteComment", DeleteComment)
		comment.GET("/getCommentCount", GetCommentCount)
		comment.POST("/listComment", ListComment)
	}
	family := engine.Group("/family")
	{
		family.POST("/saveFamily", SaveFamily)
		family.GET("/deleteFamily", DeleteFamily)
		family.GET("/getFamily", GetFamily)
		family.GET("/getAdminFamily", GetAdminFamily)
		family.POST("/listRandomFamily", ListRandomFamily)
		family.POST("/listFamily", ListFamily)
		family.GET("/changeLoveStatus", ChangeLoveStatus)
	}
	qiniu := engine.Group("/qiniu")
	{
		qiniu.GET("/getUpToken", GetUpToken)
	}
	resource := engine.Group("/resource")
	{
		resource.POST("/saveResource", SaveResource)
		resource.POST("/deleteResource", DeleteResource)
		resource.GET("/getResourceInfo", GetResourceInfo)
		resource.GET("/getImageList", GetImageList)
		resource.POST("/listResource", ListResource)
		resource.GET("/changeResourceStatus", ChangeResourceStatus)
	}
	user := engine.Group("/user")
	{
		user.POST("/regist", Register)
		user.POST("/login", Login)
		user.POST("/token", LoginByToken)
		user.GET("/logout", Logout)
		user.POST("/updateUserInfo", UpdateUserInfo)
		user.GET("/getCode", GetCaptcha)
		user.GET("/getCodeForBind", GetCaptchaForBind)
		user.POST("/updateSecretInfo", UpdateSecretInfo)
		user.GET("/getCodeForForgetPassword", GetCaptchaForForgetPassword)
		user.POST("/updateForForgetPassword", UpdateForForgetPassword)
		user.GET("/getUserByUsername", GetUserByUsername)
	}
	weiYan := engine.Group("/weiYan")
	{
		weiYan.POST("/saveWeiYan", SaveWeiYan)
		weiYan.POST("/saveNews", SaveNews)
		weiYan.POST("/listNews", ListNews)
		weiYan.GET("/deleteWeiYan", DeleteWeiYan)
		weiYan.POST("/listWeiYan", ListWeiYan)
	}
}

func Run(port int) {
	//go ws.run()

	//go func() {
	//	err := prometheus.StartPromeSrv(prometheusPort)
	//	if err != nil {
	//		panic(err)
	//	}
	//}()
}

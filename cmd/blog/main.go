package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	_ "wan_go/cmd/blog/docs"
	"wan_go/core/logger"
	"wan_go/internal/blog/router"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/log"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("logs/blog.log")

	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.Default()

	router.Init(r)

	logger.WithOutput(f)
	port := config.Config.Blog.Port[0]
	logger.Infof("localhost:%d", port)
	logger.Info("------------------------------------------------------")
	logger.Fatal(r.Run(fmt.Sprintf(":%d", port)))

}

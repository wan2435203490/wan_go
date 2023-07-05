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

	//logger.Fatal(r.Run(":8081"))
	port := config.Config.Blog.Port[0]
	logger.Fatal(r.Run(fmt.Sprintf(":%d", port)))

	logger.Info("------------------------------------------------------")

	//defaultPorts := config.Config.Blog.Port
	//ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10004 as port")
	//flag.Parse()
	//address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	//if config.Config.Api.ListenIP != "" {
	//	address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	//}
	//address = config.Config.CmsApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	//fmt.Println("start blog apis server address: ", address, ", wan_go version: ", constant.CurrentVersion, "\n")
	//
	//err := r.Run(address)
	//if err != nil {
	//	logs.Error("", "run failed ", *ginPort, err.Error())
	//}
}

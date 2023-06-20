package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	//"logs"
	"os"
	"wan_go/core/logger"
	"wan_go/internal/blog"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/log"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.ReleaseMode)
	f, err := os.Create("logs/blog.log")
	fmt.Println(err)
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.Default()

	log.NewInfo("111111", 90909)
	blog.Init(r)

	logger.Fatal(r.Run(":8081"))

	//defaultPorts := config.Config.Blog.Port
	//ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10004 as port")
	//flag.Parse()
	//address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	//if config.Config.Api.ListenIP != "" {
	//	address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	//}
	//address = config.Config.CmsApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	//fmt.Println("start blog api server address: ", address, ", wan_go version: ", constant.CurrentVersion, "\n")
	//
	//err := r.Run(address)
	//if err != nil {
	//	logs.Error("", "run failed ", *ginPort, err.Error())
	//}
}

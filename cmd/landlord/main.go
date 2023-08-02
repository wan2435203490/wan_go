package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"wan_go/core/logger"
	"wan_go/internal/landlord/router"
	"wan_go/internal/landlord/ws"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/log"
)

func main() {

	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("logs/landlord.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()

	router.Init(r)

	//ws
	go ws.Start()

	//log.Fatal(engine.RunTLS(":8080", config.Config.TLS.Cert, config.Config.TLS.Key))
	logger.WithOutput(f)
	logger.Info("------------------------------------------------------")

	port := config.Config.Landlords.Port[0]
	logger.Infof("localhost:%d", port)
	logger.Fatal(r.Run(fmt.Sprintf(":%d", port)))

}

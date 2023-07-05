package router

import (
	"github.com/gin-gonic/gin"
	"hash/crc32"
	"wan_go/internal/blog/apis"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, musicRouter)
}

func musicRouter(v1 *gin.RouterGroup) {

	api := apis.New(crc32.ChecksumIEEE)
	r := v1.Group("/music")
	{
		r.GET("/song/url", api.SongUrl)
		r.GET("/song/randomMusic", api.RandomMusic)
	}
}

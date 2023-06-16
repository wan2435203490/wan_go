package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WithCors(c *gin.Context) {

	//c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("origin"))
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*,POST, GET, PATCH, DELETE, PUT, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "*, token, content-type")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
	c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
	c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                   //  跨域请求是否需要带cookie信息 默认设置为true
	c.Header("content-type", "application/json")                                                                                                                                                           // 设置返回格式是json
	c.Header("Set-Cookie", "SameSite=None;Secure")
	//Release all option pre-requests
	if c.Request.Method == http.MethodOptions {
		fmt.Println("OPTIONS", c.Request.RequestURI)
		c.AbortWithStatusJSON(http.StatusNoContent, "Options Request!")
	}

}

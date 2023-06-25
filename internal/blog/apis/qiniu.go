package apis

import (
	"github.com/gin-gonic/gin"
	"os"
	"wan_go/pkg/common/api"
	"wan_go/pkg/utils"
)

type QiniuApi struct {
	api.Api
}

// GetUpToken 获取覆盖凭证
func (a QiniuApi) GetUpToken(c *gin.Context) {
	a.MakeContext(c)
	var key string
	//a.StringFailed(&key, "key")

	token := utils.GetQiniuToken(key)

	a.OK(token)
}

// curl --location --request GET 'http://localhost:8081/qiniu/testUpload?path=qwq.txt' \
// --header 'Content-Type: application/json'
func (a QiniuApi) TestUpload(c *gin.Context) {
	a.MakeContext(c)

	path := c.Query("path")
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	utils.UploadToQiNiu(file)
}

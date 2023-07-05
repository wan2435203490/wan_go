package actions

import (
	"net/http"
	log "wan_go/core/logger"
	"wan_go/pkg/common/dto"
	"wan_go/pkg/common/models"
	"wan_go/sdk/pkg"
	"wan_go/sdk/pkg/jwtauth/user"
	"wan_go/sdk/pkg/response"

	"github.com/gin-gonic/gin"
)

// UpdateAction 通用更新动作
func UpdateAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		req := control.Generate()
		//更新操作
		err = req.Bind(c)
		if err != nil {
			response.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
			return
		}
		var object models.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			response.Error(c, 500, err, "模型生成失败")
			return
		}
		object.SetUpdateBy(user.GetUserId(c))

		//数据权限检查
		p := GetPermissionFromContext(c)

		db = db.WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).Updates(object)
		if err = db.Error; err != nil {
			log.Errorf("MsgID[%s] Update error: %s", msgID, err)
			response.Error(c, 500, err, "更新失败")
			return
		}
		if db.RowsAffected == 0 {
			response.Error(c, http.StatusForbidden, nil, "无权更新该数据")
			return
		}
		response.OK(c, object.GetId(), "更新成功")
		c.Next()
	}
}

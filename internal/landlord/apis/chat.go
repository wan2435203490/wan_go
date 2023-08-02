package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"sync"
	msg "wan_go/internal/landlord/model/chat_msg"
	"wan_go/internal/landlord/service"
	"wan_go/internal/landlord/service/dto"
	"wan_go/internal/landlord/ws"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/sdk/pkg/jwtauth/user"
)

type ChatApi struct {
	api.Api
}

var rateLimiterMap sync.Map

func (a ChatApi) Chat(c *gin.Context) {
	s := service.Room{}
	req := dto.Chat{}
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	user := user.GetUser(c)
	if !CheckLimit(user) {
		a.ErrorInternal("你说话太快啦～")
		return
	}

	chatMsg := msg.NewChat(&req, user)

	dimensionType := landlord_const.ToDimensionType(req.Dimension)

	switch dimensionType {
	case landlord_const.Room:
		room, err := s.GetUserRoom(user.ID)
		if a.IsError(err) {
			return
		}
		nc := &ws.NotifyComponent{}
		err = nc.Send2Room(room.ID, chatMsg)
		if a.IsError(err) {
			return
		}
		a.OK()
	case landlord_const.All:
		nc := &ws.NotifyComponent{}
		err := nc.Send2AllUser(chatMsg)
		if a.IsError(err) {
			return
		}
		a.OK()
	default:
		a.ErrorInternal(fmt.Sprintf("不支持的聊天范围:%s", chatMsg.Dimension))
	}
}

// CheckLimit 限制短时间消息刷屏
func CheckLimit(user *blog.User) bool {
	limiter, ok := rateLimiterMap.Load(user.ID)
	if !ok {
		//每秒往桶中放r个令牌, 桶的容量b
		limiter = rate.NewLimiter(1, 4)
		rateLimiterMap.Store(user.ID, limiter)
	}

	return limiter.(*rate.Limiter).Allow()
}

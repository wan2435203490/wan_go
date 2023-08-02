package apis

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/landlord/service"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/utils"
)

type AchievementApi struct {
	api.Api
}

func (a AchievementApi) GetAchievementByUserId(c *gin.Context) {

	s := service.Achievement{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	userIdStr := a.Param("userId")
	if userIdStr == "" {
		a.ErrorInternal("userId is empty")
		return
	}
	userId := utils.StringToInt32(userIdStr)
	//userId := user.GetUserId32(c)

	if a.IsError(s.ExistUser(userId)) {
		return
	}

	achievement := &blog.Achievement{UserId: userId}
	if a.IsError(s.FindAchievementByUserId(achievement)) {
		return
	}

	if achievement.ID <= 0 {
		//achievement不存在
		//achievement.ID = strings.ReplaceAll(uuid.NewString(), "-", "")
		if a.IsError(s.CreateAchievement(achievement)) {
			return
		}
	}

	a.OK(achievement)
}

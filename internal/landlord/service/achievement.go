package service

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"wan_go/internal/landlord/model"
	msg "wan_go/internal/landlord/model/chat_msg"
	"wan_go/internal/landlord/service/dto"
	"wan_go/internal/landlord/ws"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/utils"
	"wan_go/sdk/service"
)

type Achievement struct {
	service.Service
}

func (s *Achievement) CountScore(user *blog.User, result *model.RoundResult) error {
	rs := &Room{}
	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		return utils.Wrap(err)
	}
	var resList []*dto.ResultScore
	var messages []*msg.GameEnd
	multiple := result.Multiple

	for _, u := range room.UserList {
		player := room.GetPlayerByUserId(u.ID)

		isWin := player.Identity == result.WinIdentity

		resultScore := dto.NewResultScore(u.UserName, multiple, isWin, player.IsLandlord())
		moneyChange, err := strconv.ParseFloat(resultScore.MoneyChange, 64)
		if err != nil {
			s.Log.Error(err)
		}
		u.Money += moneyChange
		resList = append(resList, resultScore)

		msg := msg.NewGameEnd(result.WinIdentity, isWin)
		messages = append(messages, msg)

		if err = s.UpdatesUser(u); err != nil {
			return utils.Wrap(err)
		}

		if err = s.updateAchievement(u.ID, isWin); err != nil {
			return utils.Wrap(err)
		}
	}

	nc := &ws.NotifyComponent{}
	for i, user := range room.UserList {
		endMsg := messages[i]
		endMsg.ResList = resList
		go func() {
			err := nc.Send2User(user.ID, endMsg)
			if err != nil {
				s.Log.Info(err.Error())
			}
		}()
	}

	room.Reset()
	if err = rs.UpdateRoom(room); err != nil {
		return utils.Wrap(err)
	}

	return nil
}

func (s *Achievement) updateAchievement(userId int32, isWinning bool) error {
	achievement := &blog.Achievement{UserId: userId}
	if err := s.FindAchievementByUserId(achievement); err != nil {
		return err
	}
	achievement.CalculateScore(isWinning)

	if achievement.ID <= 0 {
		//achievement.ID = strings.ReplaceAll(uuid.NewString(), "-", "")
		return s.CreateAchievement(achievement)
	} else {
		return s.UpdateAchievement(achievement)
	}
}

// gorm 暂时没必要抽象gorm层 目前AddError都在gorm层
func (s *Achievement) UpdateAchievement(achievement *blog.Achievement) error {
	return s.Orm.Debug().Updates(achievement).Error
}

func (s *Achievement) ExistUser(userId int32) error {
	err := s.Orm.Debug().Find(&blog.User{}, "id=?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}
	return nil
}

func (s *Achievement) UpdatesUser(user *blog.User) error {
	return s.Orm.Debug().Updates(user).Error
}

func (s *Achievement) FindAchievementByUserId(achievement *blog.Achievement) error {
	err := s.Orm.Debug().Where("user_id=?", achievement.UserId).Find(&achievement).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (s *Achievement) CreateAchievement(achievement *blog.Achievement) error {
	return s.Orm.Debug().Create(achievement).Error
}

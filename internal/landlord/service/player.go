package service

import (
	"wan_go/internal/landlord/model"
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/sdk/service"
)

type Player struct {
	service.Service
}

func (s *Player) GetPlayerCards(user *blog.User) ([]*model.Card, error) {
	rs := &Room{}
	return rs.GetUserCards(user.ID)
}

func (s *Player) IsPlayerRound(user *blog.User) bool {
	rs := &Room{}
	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		s.Log.Error(err.Error())
		return false
	}
	if room.RoomStatus != landlord_const.Playing {
		//todo 处理错误时的情况
		s.Log.Errorf("玩家当前状态：%s，不是Playing", room.RoomStatus)
		return false
	}
	if room.StepNum == 0 {
		//叫牌未结束
		return false
	}
	player := room.GetPlayerByUserId(user.ID)
	remain := room.StepNum % 3
	if remain == 0 {
		remain = 3
	}

	return player.ID == remain
}

func (s *Player) IsPlayerReady(user *blog.User) bool {
	rs := &Room{}
	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		s.Log.Errorf(err.Error())
		return false
	}

	if room.RoomStatus == landlord_const.Playing {
		s.Log.Errorf("房间%d游戏已经开始", room.ID)
		return false
	}
	player := room.GetPlayerByUserId(user.ID)
	return player.Ready
}

func (s *Player) CanPass(user *blog.User) bool {
	rs := &Room{}
	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		s.Log.Errorf(err.Error())
		return false
	}
	if room.RoomStatus != landlord_const.Playing {
		s.Log.Errorf("房间%d游戏尚未开始", room.ID)
		return false
	}
	player := room.GetPlayerByUserId(user.ID)
	if room.PrePlayerId == 0 {
		return player.Identity != landlord_const.Landlord
	}
	return room.PrePlayerId != player.ID
}

func (s *Player) CanBid(user *blog.User) int {
	rs := &Room{}
	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		s.Log.Errorf(err.Error())
		return -1
	}
	if room.RoomStatus != landlord_const.Playing {
		s.Log.Errorf("房间%d游戏尚未开始", room.ID)
		return -1
	}
	if room.StepNum != 0 {
		return -1
	}
	player := room.GetPlayerByUserId(user.ID)
	if player.ID != room.BiddingPlayerId {
		return -1
	}
	return room.Multiple
}

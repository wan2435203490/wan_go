package service

import (
	"errors"
	"fmt"
	"time"
	"wan_go/internal/landlord/model"
	msg "wan_go/internal/landlord/model/chat_msg"
	"wan_go/internal/landlord/service/dto"
	"wan_go/internal/landlord/ws"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/redis"
	rocksCache "wan_go/pkg/common/db/rocks_cache"
	"wan_go/pkg/utils"
	"wan_go/sdk/service"
)

type Room struct {
	service.Service
}

func (s *Room) GetRoomOut(user *blog.User, roomId int32) (*dto.RoomOut, error) {

	room, err := s.GetRoom(roomId)
	if err != nil {
		return nil, utils.Wrap(err)
	}
	if !s.canVisit(user, room) {
		return nil, errors.New("你无权查看本房间的信息")
	}
	result := dto.ToRoomOut(room)
	for _, player := range result.PlayerList {
		player.Online = ws.WS.IsOnline(player.User.ID)
	}
	s.setCountDown(room, result)
	return result, nil
}

func (s *Room) JoinRoom(user *blog.User, dtoRoom *dto.Room) error {
	dtoRoomId := dtoRoom.ID
	room, err := s.GetRoom(dtoRoomId)
	if err != nil {
		return err
	}
	if room.RoomStatus == landlord_const.Playing {
		return errors.New("房间正在游戏中，无法加入")
	}
	password := dtoRoom.Password
	if isSelf, err := s.joinRoom(dtoRoomId, password, user); isSelf || err != nil {
		return err
	}
	nc := ws.NotifyComponent{}
	err = nc.Send2Room(dtoRoomId, msg.NewPlayerJoin(user))
	//todo
	return nil
}

func (s *Room) setCountDown(room *model.Room, result *dto.RoomOut) {
	if room.PrePlayTime == 0 {
		result.CountDown = -1
		return
	}
	gap := (time.Now().UnixMilli() - room.PrePlayTime) / 1000
	countDown := config.Config.Landlords.MaxSecondsForEveryRound - gap
	if countDown <= 0 {
		result.CountDown = 0
	} else {
		result.CountDown = int(countDown)
	}
}

func (s *Room) canVisit(user *blog.User, room *model.Room) bool {
	for _, u := range room.UserList {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}

// CreateRoom 创建房间
func (s *Room) CreateRoom(user *blog.User, title, roomPassword string) (room *model.Room, err error) {
	roomId, err := rocksCache.GetUserRoomId(user.ID)
	//if err != nil {
	//	return nil, utils.Wrap(err)
	//}
	if roomId != 0 {
		return nil, errors.New(fmt.Sprintf("用户已在房间号为 %d 的房间", roomId))
	}

	roomId, err = redis.GenRoomId()
	if err != nil {
		return nil, utils.WrapMsg(err, "生成RoomId失败")
	}

	room = &model.Room{
		ID:              roomId,
		RoomStatus:      landlord_const.Preparing,
		Multiple:        0,
		PrePlayerId:     0,
		StepNum:         0,
		BiddingPlayerId: 0,
		Title:           title,
		Owner:           user,
	}

	player := &model.Player{
		ID:   1,
		User: user,
	}

	room.UserList = append(room.UserList, user)
	room.PlayerList = append(room.PlayerList, player)

	if roomPassword != "" {
		room.Locked = true
		room.Password = roomPassword
	}

	err = redis.CacheUserRoom(user.ID, roomId)
	if err != nil {
		return nil, utils.WrapMsg(err, "缓存user-room失败")
	}
	err = redis.CacheRoom(roomId, room)
	if err != nil {
		return nil, utils.WrapMsg(err, "缓存room失败")
	}

	return room, nil
}

// JoinRoom 加入房间
func (s *Room) joinRoom(roomId int32, roomPassword string, user *blog.User) (isSelf bool, err error) {
	guessRoomId, _ := rocksCache.GetUserRoomId(user.ID)
	isSelf = guessRoomId == roomId
	if isSelf {
		return
	}
	if guessRoomId != 0 {
		return false, errors.New(fmt.Sprintf("用户已在房间号为 %d 的房间", roomId))
	}

	room, err := rocksCache.GetRoom(roomId)
	if err != nil {
		return false, utils.WrapMsg(err, fmt.Sprintf("该房间不存在，房间ID: %d", roomId))
	}

	if room.ContainsUser(user) {
		//return errors.New("您已经加入此房间，无法重复加入！")
		return false, nil
	}
	if room.IsFull() {
		return false, errors.New("该房间已满，请寻找其他房间！")
	}
	if !room.CheckPassword(roomPassword) {
		return false, errors.New("对不起，您输入的房间密码有误！")
	}

	//playerId(座位序号) 可能是1 2 3 这里取不存在的player最小id，以后实现选座位
	playerId := room.GetAvailablePlayerId()
	player := &model.Player{ID: playerId, User: user}
	room.UserList = append(room.UserList, user)
	room.PlayerList = append(room.PlayerList, player)

	err = redis.CacheUserRoom(user.ID, roomId)
	if err != nil {
		return false, utils.WrapMsg(err, "缓存user-room失败")
	}

	err = redis.CacheRoom(roomId, room)
	if err != nil {
		return false, utils.WrapMsg(err, "缓存room失败")
	}

	return false, nil
}

// ExitRoom 退出房间 return 房间是否已解散
func (s *Room) ExitRoom(user *blog.User) error {
	room, err := s.GetUserRoom(user.ID)
	if err != nil {
		return err
	}

	err = rocksCache.DeleteUserRoom(user.ID)
	if err != nil {
		return err
	}

	room.RemoveUser(user.ID)
	room.RemovePlayer(user.ID)

	if len(room.PlayerList) == 0 {
		err = rocksCache.DeleteRoom(room.ID)
		if err != nil {
			return err
		}
		return nil
	}

	err = redis.CacheRoom(room.ID, room)
	if err != nil {
		return err
	}

	nc := ws.NotifyComponent{}
	err = nc.Send2Room(room.ID, msg.NewPlayerExit(user))
	//todo
	//if err != nil {
	//	return err
	//}

	return nil
}

func (s *Room) ListRooms() (*[]*model.Room, error) {
	return redis.ListRoom()
}

func (s *Room) GetRoom(roomId int32) (*model.Room, error) {
	room, err := rocksCache.GetRoom(roomId)
	if err != nil {
		return nil, utils.WrapMsg(err, "该房间不存在，请核实您输入的房间号！")
	}
	return room, nil
}

func (s *Room) UpdateRoom(new *model.Room) error {
	if !redis.ExistsRoom(new.ID) {
		return errors.New("该房间不存在，请刷新界面重试！")
	}

	return redis.CacheRoom(new.ID, new)
}

// GetUserRoom 获取当前用户所在的房间对象
func (s *Room) GetUserRoom(userId int32) (*model.Room, error) {

	roomId, err := rocksCache.GetUserRoomId(userId)
	if err != nil {
		return nil, utils.WrapMsg(err, fmt.Sprintf("当前用户（%d）不在任何房间", userId))
	}

	room, err := rocksCache.GetRoom(roomId)
	if err != nil {
		return nil, utils.WrapMsg(err, fmt.Sprintf("未找到对应房间:%d", roomId))
	}

	return room, nil
}

func (s *Room) GetUserCards(userId int32) ([]*model.Card, error) {
	room, err := s.GetUserRoom(userId)
	if err != nil {
		return nil, err
	}

	player := room.GetPlayerByUserId(userId)
	if player == nil {
		return nil, errors.New("未找到该玩家！")
	}
	return player.Cards, nil
}

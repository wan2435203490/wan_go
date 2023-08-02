package msg

import (
	"wan_go/pkg/common/constant/landlord_const"
)

type startGame struct {
	Message
	RoomId int32 `json:"roomId"`
}

func NewStartGame(roomId int32) *startGame {
	var v startGame
	v.Type = v.GetMessageType()
	v.RoomId = roomId
	return &v
}

func (s *startGame) GetMessageType() string {
	return landlord_const.StartGame.GetWsMessageType()
}

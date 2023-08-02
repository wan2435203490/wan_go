package msg

import "wan_go/pkg/common/constant/landlord_const"

type unReadyGame struct {
	Message
	UserId int32 `json:"userId"`
}

func NewUnReadyGame(userId int32) *unReadyGame {
	var v unReadyGame
	v.UserId = userId
	v.Type = v.GetMessageType()
	return &v
}

func (u *unReadyGame) GetMessageType() string {
	return landlord_const.UnReadyGame.GetWsMessageType()
}

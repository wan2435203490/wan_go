package msg

import (
	"wan_go/pkg/common/constant/landlord_const"
)

type readyGame struct {
	Message
	UserId int32 `json:"userId"`
}

func NewReadyGame(userId int32) *readyGame {
	var v readyGame
	v.UserId = userId
	v.Type = v.GetMessageType()
	return &v
}

func (r *readyGame) GetMessageType() string {
	return landlord_const.ReadyGame.GetWsMessageType()
}

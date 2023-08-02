package msg

import (
	"wan_go/internal/landlord/service/dto"
	"wan_go/pkg/common/constant/landlord_const"
)

type GameEnd struct {
	Message
	WinningIdentity landlord_const.Identity `json:"winingIdentity"`
	IsWinning       bool                    `json:"winning"`
	ResList         []*dto.ResultScore      `json:"resList"`
}

func EmptyGameEnd() *GameEnd {
	v := &GameEnd{}
	v.Type = v.GetMessageType()
	return v
}

func NewGameEnd(winningIdentity landlord_const.Identity, isWinning bool) *GameEnd {
	v := &GameEnd{
		WinningIdentity: winningIdentity,
		IsWinning:       isWinning,
	}
	v.Type = v.GetMessageType()
	return v
}

//
//func NewGameEnd(winningIdentity landlord_const.Identity, isWinning bool,
//	resList []*DTO.ResultScore) *GameEnd {
//	v := &GameEnd{
//		WinningIdentity: winningIdentity,
//		IsWinning:       isWinning,
//		ResList:         resList,
//	}
//	v.Type = v.GetMessageType()
//	return v
//}

func (g *GameEnd) GetMessageType() string {
	return landlord_const.GameEnd.GetWsMessageType()
}

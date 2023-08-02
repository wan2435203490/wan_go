package dto

import (
	"wan_go/internal/landlord/model"
	"wan_go/pkg/common/constant/landlord_const"
)

type PlayerOut struct {
	ID          int                     `json:"id"`
	CardSize    int                     `json:"cardSize"`
	Identity    landlord_const.Identity `json:"identity"`
	RecentCards []*model.Card           `json:"recentCards"`
	Ready       bool                    `json:"ready"`
	Online      bool                    `json:"online"`
	User        *UserOut                `json:"user"`
}

func (p *PlayerOut) GetIdentityName() string {
	//todo 如何表示空？ 取一个default值？
	return p.Identity.GetIdentity()
}

func ToPlayerOut(p *model.Player) *PlayerOut {
	if p == nil {
		return nil
	}

	return &PlayerOut{
		ID:          p.ID,
		CardSize:    len(p.Cards),
		Identity:    p.Identity,
		Ready:       p.Ready,
		User:        ToUserOut(p.User),
		RecentCards: p.RecentCards,
	}
}

func ToPlayerOutList(players []*model.Player) []*PlayerOut {
	if players == nil {
		return nil
	}

	var ret []*PlayerOut

	for _, player := range players {
		ret = append(ret, ToPlayerOut(player))
	}

	return ret
}

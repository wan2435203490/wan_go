package msg

import (
	"wan_go/pkg/common/constant/landlord_const"
)

type pleasePlayCard struct {
	Message
}

func NewPleasePlayCard() *pleasePlayCard {
	var v pleasePlayCard
	v.Type = v.GetMessageType()
	return &v
}

func (p *pleasePlayCard) GetMessageType() string {
	return landlord_const.PleasePlayCard.GetWsMessageType()
}

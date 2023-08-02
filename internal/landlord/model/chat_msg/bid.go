package msg

import "wan_go/pkg/common/constant/landlord_const"

type Bid struct {
	Message
	Score int `json:"score"`
}

func NewBid(score int) *Bid {
	var v Bid
	v.Type = v.GetMessageType()
	v.Score = score
	return &v
}

func (b *Bid) GetMessageType() string {
	return landlord_const.Bid.GetWsMessageType()
}

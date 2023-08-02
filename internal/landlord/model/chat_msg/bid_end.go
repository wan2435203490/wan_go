package msg

import "wan_go/pkg/common/constant/landlord_const"

type BidEnd struct {
	Message
}

func NewBidEnd() *BidEnd {
	var v BidEnd
	v.Type = v.GetMessageType()
	return &v
}

func (b *BidEnd) GetMessageType() string {
	return landlord_const.BidEnd.GetWsMessageType()
}

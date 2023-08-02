package msg

import (
	"time"
	"wan_go/pkg/common/constant/landlord_const"
)

type pong struct {
	Message
	TimeStamp time.Time `json:"timeStamp"`
}

func NewPong() *pong {
	var v pong
	v.TimeStamp = time.Now()
	v.Type = v.GetMessageType()
	return &v
}

func (p *pong) GetMessageType() string {
	return landlord_const.Pong.GetWsMessageType()
}

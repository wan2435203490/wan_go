package dto

import (
	"wan_go/internal/landlord/model"
	"wan_go/pkg/common/constant/landlord_const"
)

type RoomOut struct {
	ID         int32                     `json:"id"`
	Title      string                    `json:"title"`
	Owner      *UserOut                  `json:"owner"`
	PlayerList []*PlayerOut              `json:"playerList"`
	RoomStatus landlord_const.RoomStatus `json:"roomStatus"`
	Status     string                    `json:"status"`
	Multiple   int                       `json:"multiple"`
	TopCards   []*model.Card             `json:"topCards"`
	StepNum    int                       `json:"stepNum"`
	CountDown  int                       `json:"countdown"`
}

func ToRoomOut(r *model.Room) *RoomOut {
	out := &RoomOut{
		ID:         r.ID,
		Title:      r.Title,
		Owner:      ToUserOut(r.Owner),
		PlayerList: ToPlayerOutList(r.PlayerList),
		RoomStatus: r.RoomStatus,
		Status:     r.RoomStatus.GetRoomStatus(),
		Multiple:   r.Multiple,
		StepNum:    r.StepNum,
	}

	if r.Distribution != nil {
		out.TopCards = r.Distribution.TopCards
	}

	return out
}

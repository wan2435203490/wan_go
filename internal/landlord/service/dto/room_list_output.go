package dto

import (
	"wan_go/internal/landlord/model"
	"wan_go/pkg/common/constant/landlord_const"
)

type RoomListOutput struct {
	ID         int32                     `json:"id"`
	Title      string                    `json:"title"`
	Owner      *UserOut                  `json:"owner"`
	Locked     bool                      `json:"locked"`
	UserList   []*UserOut                `json:"userList"`
	RoomStatus landlord_const.RoomStatus `json:"roomStatus"`
}

func RoomListOutputFromRoom(room *model.Room) *RoomListOutput {
	return &RoomListOutput{
		ID:         room.ID,
		Title:      room.Title,
		Owner:      ToUserOut(room.Owner),
		Locked:     room.Locked,
		UserList:   ToUserOutList(room.UserList),
		RoomStatus: room.RoomStatus,
	}
}

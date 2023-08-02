package landlord_const

import (
	"github.com/bytedance/sonic"
)

type RoomStatus int

const (
	Preparing RoomStatus = iota
	Playing
)

func (s RoomStatus) GetRoomStatusName() string {
	return []string{"准备中", "游戏中"}[s]
}

func (s RoomStatus) GetRoomStatus() string {
	return []string{"PREPARING", "PLAYING"}[s]
}

func (s RoomStatus) MarshalJSON() ([]byte, error) {
	var str string
	switch s {
	case Preparing:
		str = "PREPARING"
	case Playing:
		str = "PLAYING"
	}
	return sonic.Marshal(str)
}

func (s *RoomStatus) UnmarshalJSON(b []byte) error {

	var str string
	if err := sonic.Unmarshal(b, &str); err != nil {
		return err
	}
	switch str {
	case "PREPARING":
		*s = Preparing
	case "PLAYING":
		*s = Playing
	}

	return nil
}

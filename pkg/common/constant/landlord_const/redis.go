package landlord_const

import "fmt"

const (
	RocksCachePrefix = "RC"
	UserRoomKey      = "UserRoom"
	RoomKey          = "Room"
	RoomIdentity     = "RoomIdentity"
)

func GetUserRoomKey(userId int32) string {
	return fmt.Sprintf("%s-%d", UserRoomKey, userId)
}

func GetRoomKey(roomId int32) string {
	return fmt.Sprintf("%s-%d", RoomKey, roomId)
}

func GetRocksCacheKey(key string) string {
	return RocksCachePrefix + key
}

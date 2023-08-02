package rocksCache

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"wan_go/internal/landlord/model"
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/redis"
	"wan_go/pkg/common/log"
	"wan_go/pkg/utils"
)

func DelLandlordKeys() {
	fmt.Println("init rocks cache to del old keys")
	fName := utils.GetSelfFuncName()
	for _, key := range []string{landlord_const.RoomKey, landlord_const.UserRoomKey} {
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			keys, cursor, err = db.DB.RDB.Scan(context.Background(), cursor, landlord_const.RocksCachePrefix+key+"*", 3000).Result()
			if err != nil {
				panic(err.Error())
			}
			n += len(keys)
			// for each for redis cluster
			for _, key := range keys {
				if err = db.DB.RDB.Del(context.Background(), key).Err(); err != nil {
					log.NewError("", fName, key, err.Error())
					err = db.DB.RDB.Del(context.Background(), key).Err()
					if err != nil {
						panic(err.Error())
					}
				}
			}
			if cursor == 0 {
				break
			}
		}
	}
}

func GetUserRoomId(userId int32) (int32, error) {
	key := landlord_const.GetUserRoomKey(userId)

	//str, err := db.DB.Rc.Fetch(landlord_const.GetRocksCacheKey(key), time.Hour*24, func() (string, error) {
	//	val := db.DB.RDB.Get(context.Background(), key)
	//	return val.Val(), val.Err()
	//})
	//roomId := utils.StringToInt32(str)
	//return roomId, utils.Wrap(err)

	val := db.DB.RDB.Get(context.Background(), key)
	return utils.StringToInt32(val.Val()), val.Err()
}

func DeleteUserRoom(userId int32) error {
	key := landlord_const.GetUserRoomKey(userId)
	var err error
	//if err = db.DB.Rc.TagAsDeleted(landlord_const.GetRocksCacheKey(key)); err != nil {
	//	return utils.Wrap(err)
	//}
	if err = redis.Delete(key); err != nil {
		return utils.Wrap(err)
	}
	return utils.Wrap(err)
}

func GetRoom(roomId int32) (*model.Room, error) {
	key := landlord_const.GetRoomKey(roomId)
	var err error
	//str, err := db.DB.Rc.Fetch(landlord_const.GetRocksCacheKey(key), time.Hour*24, func() (string, error) {
	//	val := db.DB.RDB.Get(context.Background(), key)
	//	return val.Val(), val.Err()
	//})
	//if err != nil {
	//	return nil, utils.Wrap(err)
	//}
	val := db.DB.RDB.Get(context.Background(), key)
	str, err := val.Val(), val.Err()
	var room model.Room
	err = sonic.UnmarshalString(str, &room)
	return &room, utils.Wrap(err)
}

func DeleteRoom(roomId int32) error {
	key := landlord_const.GetRoomKey(roomId)
	var err error
	//if err = db.DB.Rc.TagAsDeleted(landlord_const.GetRocksCacheKey(key)); err != nil {
	//	return utils.Wrap(err)
	//}
	if err = redis.Delete(key); err != nil {
		return utils.Wrap(err)
	}
	return nil
}

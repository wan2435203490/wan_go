package redis

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
	"time"
	"wan_go/internal/landlord/model"
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/utils"
)

func CacheUserRoom(userId, roomId int32) error {
	key := landlord_const.GetUserRoomKey(userId)
	//rdb().Expire(context.Background(), key, time.Hour*24)
	//todo 房间自定义设置时长
	return rdb().Set(context.Background(), key, roomId, time.Hour*24).Err()
}

//
//func ExistsUserRoom(userId int32) bool {
//	key := landlord_const.GetUserRoomKey(userId)
//	return rdb().Exists(context.Background(), key).Val() != 0
//}

func CacheRoom(roomId int32, room *model.Room) error {
	key := landlord_const.GetRoomKey(roomId)
	marshal, err := sonic.Marshal(room)
	if err != nil {
		return utils.WrapMsg(err, "序列化错误")
	}
	//todo 房间自定义设置时长
	//rdb().Expire(context.Background(), key, time.Hour*24)
	return rdb().Set(context.Background(), key, marshal, time.Hour*24).Err()
}

func DeleteRoom(roomId int32) error {
	key := landlord_const.GetRoomKey(roomId)
	return rdb().Del(context.Background(), key).Err()
}

func Delete(key string) error {
	return rdb().Del(context.Background(), key).Err()
}

func ListRoom() (rooms *[]*model.Room, err error) {

	rr := make([]*model.Room, 0)

	key := landlord_const.RoomKey
	ctx := context.Background()
	iter := rdb().Scan(ctx, 0, key+"*", 100).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		d, err := rdb().TTL(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		if d == -1 { // -1 means no TTL
			if err = rdb().Del(ctx, key).Err(); err != nil {
				return nil, err
			}
		}

		val := rdb().Get(ctx, key).Val()
		room := model.Room{}
		err = sonic.UnmarshalString(val, &room)
		if err != nil {
			return nil, err
		}
		rr = append(rr, &room)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return &rr, nil
}

// GenRoomId 生成自增id
func GenRoomId() (int32, error) {

	genIdScript := redis.NewScript(`
local id_key = KEYS[1]
local current = redis.call('get',id_key)
if current == false then
	redis.call('set',id_key,1)
	return '10000'
end
--redis.log(redis.LOG_NOTICE,' current:'..current..':')
local result = tonumber(current)+1
--redis.log(redis.LOG_NOTICE,' result:'..result..':')
redis.call('set',id_key,result)
return result
`)

	n, err := genIdScript.Run(context.Background(), rdb(), []string{landlord_const.RoomIdentity}, 2).Result()
	if err != nil {
		return -1, err
	} else {
		//var roomId = int32(n.(string))
		var roomId = utils.StringToInt32(n.(string))
		return roomId, nil
		//id, err := strconv.Atoi(str)
		//if err == nil {
		//	return id, err
		//} else {
		//	return -1, err
		//}
	}
}

func ExistsRoom(roomId int32) bool {
	key := landlord_const.GetRoomKey(roomId)
	return rdb().Exists(context.Background(), key).Val() > 0
}

//
//func UpdateRoom(room *model.Room) error {
//	key := landlord_const.GetRoomKey(room.ID)
//	rdb().Set(context.Background(), key, room)
//}

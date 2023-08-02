package ws

import (
	"encoding/json"
	"fmt"
	"wan_go/core/logger"
	"wan_go/internal/landlord/model"
	ws "wan_go/internal/landlord/model/chat_msg"
	rocksCache "wan_go/pkg/common/db/rocks_cache"
	"wan_go/pkg/utils"
)

type NotifyComponent struct {
}

func GetRoom(roomId int32) (*model.Room, error) {
	room, err := rocksCache.GetRoom(roomId)
	if err != nil {
		return nil, utils.WrapMsg(err, "该房间不存在，请核实您输入的房间号！")
	}
	return room, nil
}

func (nc *NotifyComponent) SendStr2Room(roomId int32, content string) error {
	room, err := GetRoom(roomId)
	if err != nil {
		return err
	}
	ids := room.GetUserIds()
	logger.Info(fmt.Sprintf("SendStr2Room:%d, %s", roomId, content))
	err = WS.Send2Users(ids, content)
	//todo
	if err != nil {
		logger.Error(err)
	}
	return nil
}

func (nc *NotifyComponent) Send2Room(roomId int32, msg ws.IMessage) error {
	if bs, err := json.Marshal(msg); err != nil {
		return err
	} else {
		logger.Info(fmt.Sprintf("Send2Room:%d, %s", roomId, string(bs)))
		return nc.SendStr2Room(roomId, string(bs))
	}
}

func (nc *NotifyComponent) Send2User(userId int32, msg ws.IMessage) error {
	if bs, err := json.Marshal(msg); err != nil {
		return err
	} else {
		logger.Info(fmt.Sprintf("Send2User:%d, %s", userId, string(bs)))
		err = WS.Send2User(userId, string(bs))
		//todo
		if err != nil {
			logger.Error(err)
		}
		return nil
	}
}

func (nc *NotifyComponent) Send2AllUser(msg ws.IMessage) error {
	if bs, err := json.Marshal(msg); err != nil {
		return err
	} else {
		logger.Info("Send2AllUser:" + string(bs))
		err = WS.Send2AllUser(string(bs))
		//todo
		if err != nil {
			logger.Error(err)
		}
		return nil
	}
}

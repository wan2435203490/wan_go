package apis

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/landlord/service"
	"wan_go/internal/landlord/service/dto"
	"wan_go/pkg/common/api"
	"wan_go/pkg/utils"
	user2 "wan_go/sdk/pkg/jwtauth/user"
)

type RoomApi struct {
	api.Api
}

func (a RoomApi) Rooms(c *gin.Context) {
	s := &service.Room{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	rooms, err := s.ListRooms()
	if a.IsError(err) {
		return
	}
	if rooms != nil && len(*rooms) != 0 {
		var roomsOut []*dto.RoomListOutput
		for _, room := range *rooms {
			roomsOut = append(roomsOut, dto.RoomListOutputFromRoom(room))
		}
		a.OK(roomsOut)
		return
	}
	a.OK(rooms)
}

func (a RoomApi) GetById(c *gin.Context) {
	s := &service.Room{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	roomId := a.Param("id")
	if roomId == "" {
		return
	}
	user := user2.GetUser(c)

	outRoom, err := s.GetRoomOut(user, utils.StringToInt32(roomId))
	if a.IsError(err) {
		return
	}

	a.OK(outRoom)
}

func (a RoomApi) Create(c *gin.Context) {
	s := &service.Room{}
	var req dto.CreateRoom
	if a.MakeContextChain(c, &s.Service, &req) == nil {
		return
	}

	if len(req.Title) == 0 {
		a.ErrorInternal("房间名称不能为空")
		return
	}
	user := user2.GetUser(c)
	room, err := s.CreateRoom(user, req.Title, req.Password)
	if a.IsError(err) {
		return
	}

	a.OK(room)
}

func (a RoomApi) Join(c *gin.Context) {
	var room dto.Room
	s := &service.Room{}
	if a.MakeContextChain(c, &s.Service, &room) == nil {
		return
	}

	user := user2.GetUser(c)
	if a.IsError(s.JoinRoom(user, &room)) {
		return
	}

	a.OK("加入成功")
}

func (a RoomApi) Exit(c *gin.Context) {
	var room dto.Room
	s := &service.Room{}
	if a.MakeContextChain(c, &s.Service, &room) == nil {
		return
	}
	user := user2.GetUser(c)

	if a.IsError(s.ExitRoom(user)) {
		return
	}
	a.OK()
}

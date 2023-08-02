package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"wan_go/internal/landlord/model"
	"wan_go/internal/landlord/service"
	"wan_go/internal/landlord/service/dto"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/db/mysql/blog"
	user2 "wan_go/sdk/pkg/jwtauth/user"
)

type GameApi struct {
	api.Api
}

func (a GameApi) Ready(c *gin.Context) {
	s := service.Game{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	user := user2.GetUser(c)
	a.OK(s.ReadyGame(user))
}

func (a GameApi) UnReady(c *gin.Context) {
	s := service.Game{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	user := user2.GetUser(c)
	a.OK(s.UnReadyGame(user))
}

func (a GameApi) Bid(c *gin.Context) {
	s := service.Game{}
	var bid dto.Bid
	if a.MakeContextChain(c, &s.Service, &bid) == nil {
		return
	}

	user := user2.GetUser(c)
	if bid.Want {
		if a.IsError(s.Want(user, bid.Score)) {
			return
		}
		a.OK("已叫地主并分配身份")
	} else {
		if a.IsError(s.NoWant(user)) {
			return
		}
		a.OK("已选择不叫地主，并传递给下家")
	}
}

func (a GameApi) Play(c *gin.Context) {
	s := service.Game{}
	var cardList []*model.Card
	if a.MakeContextChain(c, &s.Service, &cardList, binding.JSON) == nil {
		return
	}

	user := user2.GetUser(c)
	if !a.validRound(user) {
		return
	}
	//if !a.Bind(&cardList, binding.JSON) {
	//	return
	//}

	//todo 逻辑优化 职责不单一
	result, err := s.PlayCard(user, cardList)
	if a.IsError(err) {
		return
	}

	as := service.Achievement{}
	a.MakeService(&as.Service)

	if result == nil {
		a.OK()
	} else {
		if a.IsError(as.CountScore(user, result)) {
			return
		}
		a.OK(result)
	}
}

func (a GameApi) Pass(c *gin.Context) {

	s := service.Game{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	user := user2.GetUser(c)

	if !a.validRound(user) {
		return
	}
	if a.IsError(s.PassGame(user)) {
		return
	}
	a.OK()
}

func (a GameApi) Give(c *gin.Context) {

	s := service.Game{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	user := user2.GetUser(c)

	if !a.validRound(user) {
		return
	}
	cards, err := s.GiveCards(user)
	if a.IsError(err) {
		return
	}
	a.OK(cards)
}

// validRound valid 需要 set response
func (a GameApi) validRound(user *blog.User) bool {

	ps := service.Player{}
	a.MakeService(&ps.Service)

	isPlayerRound := ps.IsPlayerRound(user)
	if !isPlayerRound {
		log.Printf("当前不是该玩家(%d)出牌回合", user.ID)
	}
	return isPlayerRound
}

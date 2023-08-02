package apis

import (
	"github.com/gin-gonic/gin"
	"sort"
	"wan_go/internal/landlord/service"
	"wan_go/pkg/common/api"
	user2 "wan_go/sdk/pkg/jwtauth/user"
)

type PlayerApi struct {
	api.Api
}

func (a PlayerApi) Cards(c *gin.Context) {
	s := service.Player{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	user := user2.GetUser(c)

	cards, err := s.GetPlayerCards(user)
	if a.IsError(err) {
		return
	}

	sort.SliceStable(cards, func(i, j int) bool {
		return cards[i].Grade > cards[j].Grade
	})
	a.OK(cards)
}

func (a PlayerApi) Round(c *gin.Context) {
	s := service.Player{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	user := user2.GetUser(c)

	round := s.IsPlayerRound(user)
	a.OK(round)
}

func (a PlayerApi) PlayerReady(c *gin.Context) {
	s := service.Player{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	user := user2.GetUser(c)
	ready := s.IsPlayerReady(user)
	a.OK(ready)
}

func (a PlayerApi) PlayerPass(c *gin.Context) {
	s := service.Player{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	user := user2.GetUser(c)
	can := s.CanPass(user)
	a.OK(can)
}

func (a PlayerApi) Bidding(c *gin.Context) {
	s := service.Player{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	user := user2.GetUser(c)
	score := s.CanBid(user)
	a.OK(score)
}

package model

import (
	"sort"
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db/mysql/blog"
)

type Player struct {
	//1,2,3
	ID          int                     `json:"id"`
	Identity    landlord_const.Identity `json:"identity"`
	Cards       []*Card                 `json:"cards"`
	RecentCards []*Card                 `json:"recentCards"`
	User        *blog.User              `json:"user"`
	Ready       bool                    `json:"ready"`
}

func (p *Player) GetNextPlayerId() int {
	//return utils.IfThen(p.ID == 3, 1, p.ID+1).(int)
	return p.ID%3 + 1
}

func (p *Player) AddCards(cards []*Card) {
	p.Cards = append(p.Cards, cards...)

	sort.SliceStable(p.Cards, func(i, j int) bool {
		return p.Cards[i].Grade > p.Cards[j].Grade
	})
	//sort.Ints(p.Cards)
}

func (p *Player) RemoveCards(cards []*Card) {
	for _, card := range cards {
		for j, old := range p.Cards {
			if old.Equals(card) {
				p.Cards = append(p.Cards[:j], p.Cards[j+1:]...)
				break
			}
		}
	}
}

func (p *Player) ClearRecentCards() {
	p.RecentCards = nil
}

func (p *Player) Reset() {
	p.Cards = nil
	p.Ready = false
	p.Identity = -1
	p.RecentCards = nil
}

func (p *Player) IsLandlord() bool {
	return p.Identity == landlord_const.Landlord
}

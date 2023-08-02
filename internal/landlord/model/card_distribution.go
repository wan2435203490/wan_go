package model

import (
	"math/rand"
	"sort"
	"time"
)

const (
	PokerCount = 54
)

// 每局开始重新发牌、分牌、获取每个玩家的牌
type CardDistribution struct {
	AllCardList  []*Card
	Player1Cards []*Card
	Player2Cards []*Card
	Player3Cards []*Card
	TopCards     []*Card
}

func (d *CardDistribution) Refresh() {
	d.Clear()
	d.CreateCards()
	d.Shuffle()
	d.Deal()
}

func (d *CardDistribution) Clear() {
	d.AllCardList = nil
	d.Player1Cards = nil
	d.Player2Cards = nil
	d.Player3Cards = nil
	d.TopCards = nil
}

func (d *CardDistribution) CreateCards() {
	for i := 0; i < PokerCount; i++ {
		d.AllCardList = append(d.AllCardList, GetCard(i))
	}
}

func (d *CardDistribution) Shuffle() {
	rand.Seed(time.Now().Unix())
	rand.Shuffle(PokerCount, func(i, j int) {
		d.AllCardList[i], d.AllCardList[j] = d.AllCardList[j], d.AllCardList[i]
	})
}

func (d *CardDistribution) Deal() {
	d.Player1Cards = append(d.Player1Cards, d.AllCardList[:17]...)
	d.Player2Cards = append(d.Player2Cards, d.AllCardList[17:34]...)
	d.Player3Cards = append(d.Player3Cards, d.AllCardList[34:51]...)
	d.TopCards = d.AllCardList[51:]

	SortCards(d.Player1Cards)
	SortCards(d.Player2Cards)
	SortCards(d.Player3Cards)
	//sort.Slice(d.Player1Cards, func(i, j int) bool {
	//	return d.Player1Cards[i].Grade < d.Player1Cards[j].Grade
	//})
	//sort.Slice(d.Player2Cards, func(i, j int) bool {
	//	return d.Player2Cards[i].Grade < d.Player2Cards[j].Grade
	//})
	//sort.Slice(d.Player3Cards, func(i, j int) bool {
	//	return d.Player3Cards[i].Grade < d.Player3Cards[j].Grade
	//})

}

func SortCards(cards []*Card) {
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Grade == cards[j].Grade {
			return cards[i].Type < cards[j].Type
		}
		return cards[i].Grade < cards[j].Grade
	})
}

func (d *CardDistribution) GetCards(number int) []*Card {
	switch number {
	case 1:
		return d.Player1Cards
	case 2:
		return d.Player2Cards
	case 3:
		return d.Player3Cards
	}
	panic("GetCardsByNumber number can only in 1,2,3")
}

package msg

import (
	"sort"
	"wan_go/internal/landlord/model"
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db/mysql/blog"
)

type playCard struct {
	Message
	User     *blog.User                `json:"user"`
	CardList []*model.Card             `json:"cardList"`
	CardType landlord_const.Type       `json:"cardType"`
	Number   landlord_const.CardNumber `json:"number"`
}

func NewPlayCard(user *blog.User, cardList []*model.Card, cardType landlord_const.Type) *playCard {
	card := &playCard{
		User:     user,
		CardType: cardType,
	}

	card.Type = card.GetMessageType()

	sort.Slice(cardList, func(i, j int) bool {
		return cardList[i].Grade > cardList[j].Grade
	})

	if cardType == landlord_const.Single || cardType == landlord_const.Pair {
		card.Number = cardList[0].Number
	}

	return card
}

func (p *playCard) GetMessageType() string {
	return landlord_const.PlayCard.GetWsMessageType()
}

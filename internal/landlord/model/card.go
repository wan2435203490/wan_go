package model

import (
	"wan_go/pkg/common/constant/landlord_const"
)

type Card struct {
	Id     int                       `json:"id"`
	Type   landlord_const.CardType   `json:"type"`
	Number landlord_const.CardNumber `json:"number"`
	Grade  landlord_const.CardGrade  `json:"grade"`
}

func (c *Card) CompareTo(c2 *Card) bool {
	return c.Grade.GreatThanGrade(c2.Grade)
}

//func GreatThanGrade(i, j int) bool {
//	return p.Cards[i].Grade > p.Cards[j].Grade
//}

func (c *Card) EqualsByGrade(c2 *Card) bool {
	return c.Grade == c2.Grade
}

func (c *Card) Equals(c2 *Card) bool {
	return c.Grade == c2.Grade && c.Type == c2.Type
}

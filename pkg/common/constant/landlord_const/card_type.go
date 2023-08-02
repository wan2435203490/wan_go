package landlord_const

import (
	"encoding/json"
	"strings"
)

type CardType int

const (
	//♠️
	Spade CardType = iota
	//♥️
	Heart
	//♣️
	Club
	//♦️
	Diamond
	//🃏小王
	SmallJokerType
	//大王
	BigJokerType
)

func (c *CardType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToUpper(s) {
	default:
		*c = Spade
	case "SPADE":
		*c = Spade
	case "HEART":
		*c = Heart
	case "CLUB":
		*c = Club
	case "DIAMOND":
		*c = Diamond
	case "SMALL_JOKER":
		*c = SmallJokerType
	case "BIG_JOKER":
		*c = BigJokerType
	}

	return nil
}

func (c CardType) MarshalJSON() ([]byte, error) {
	var s string
	switch c {
	default:
		s = "SPADE"
	case Spade:
		s = "SPADE"
	case Heart:
		s = "HEART"
	case Club:
		s = "CLUB"
	case Diamond:
		s = "DIAMOND"
	case SmallJokerType:
		s = "SMALL_JOKER"
	case BigJokerType:
		s = "BIG_JOKER"
	}

	return json.Marshal(s)
}

func (c CardType) GetCardType() string {
	return []string{"SPADE", "HEART", "CLUB", "DIAMOND", "SMALL_JOKER", "BIG_JOKER"}[c]
}

func (c CardType) GetCardTypeName() string {
	return []string{"黑桃", "红桃", "梅花", "方块", "小王", "大王"}[c]
}

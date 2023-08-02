package util

import (
	"wan_go/internal/landlord/model"
	"wan_go/pkg/common/constant/landlord_const"
)

func CanPlayCards(cards, preCards []*model.Card, typ, preType landlord_const.Type) bool {
	if cards == nil || preCards == nil ||
		typ == -1 || preType == -1 ||
		preType == landlord_const.JokerBomb {
		return false
	}

	if typ == landlord_const.JokerBomb ||
		preType != landlord_const.Bomb && typ == landlord_const.Bomb {
		return true
	}

	if typ == landlord_const.Bomb && preType == landlord_const.Bomb {
		return cards[0].CompareTo(preCards[0])
	}

	//当前出的不是炸弹 牌型不一样不能出
	if preType != typ {
		return false
	}

	switch typ {
	case landlord_const.ThreeWithOne:
		return cards[1].CompareTo(preCards[1])
	case landlord_const.ThreeWithPair:
		return cards[2].CompareTo(preCards[2])
	case landlord_const.FourWithTwo:
		return cards[3].CompareTo(preCards[3])
	case landlord_const.FourWithFour:
		_, i0 := IsFourWithFour(cards)
		_, i1 := IsFourWithFour(preCards)
		return i0.GreatThanGrade(i1)
	case landlord_const.Straight:
		fallthrough
	case landlord_const.StraightPair:
		if len(cards) != len(preCards) {
			return false
		}
		fallthrough
	case landlord_const.Single:
		fallthrough
	case landlord_const.Pair:
		fallthrough
	case landlord_const.ThreeType:
		fallthrough
	case landlord_const.Aircraft:
		return cards[0].CompareTo(preCards[0])
	case landlord_const.AircraftWithSingleWings:
		_, i0 := IsAircraftWithSingleWing(cards)
		_, i1 := IsAircraftWithSingleWing(preCards)
		return i0.GreatThanGrade(i1)
	case landlord_const.AircraftWithPairWings:
		_, i0 := IsAircraftWithPairWing(cards)
		_, i1 := IsAircraftWithPairWing(preCards)
		return i0.GreatThanGrade(i1)
	}

	return false
}

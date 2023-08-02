package util

import (
	"sort"
	"wan_go/internal/landlord/model"
	"wan_go/pkg/common/constant/landlord_const"
)

// 抽象type 获取一个可比较的值出来
func GetCardsType(cards ...*model.Card) landlord_const.Type {
	if IsSingle(cards) {
		return landlord_const.Single
	}
	if IsPair(cards) {
		return landlord_const.Pair
	}
	if IsThree(cards) {
		return landlord_const.ThreeType
	}
	if IsThreeWithOne(cards) {
		return landlord_const.ThreeWithOne
	}
	if IsThreeWithPair(cards) {
		return landlord_const.ThreeWithPair
	}
	if IsStraight(cards, 5) {
		return landlord_const.Straight
	}
	if IsStraightPair(cards) {
		return landlord_const.StraightPair
	}
	if IsFourWithTwo(cards) {
		return landlord_const.FourWithTwo
	}
	if b, _ := IsFourWithFour(cards); b {
		return landlord_const.FourWithFour
	}
	if IsBomb(cards) {
		return landlord_const.Bomb
	}
	if IsJokerBomb(cards) {
		return landlord_const.JokerBomb
	}
	if IsAircraft(cards) {
		return landlord_const.Aircraft
	}
	if b, _ := IsAircraftWithSingleWing(cards); b {
		return landlord_const.AircraftWithSingleWings
	}
	if b, _ := IsAircraftWithPairWing(cards); b {
		return landlord_const.AircraftWithPairWings
	}
	return -1

}

// 对牌进行从小到大地排序 花色？
func SortCards(cards []*model.Card) {
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Grade == cards[j].Grade {
			return cards[i].Type < cards[j].Type
		}
		return cards[i].Grade < cards[j].Grade
	})
}

// 从大到小
func SortCardsDesc(cards []*model.Card) {
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Grade == cards[j].Grade {
			return cards[i].Type > cards[j].Type
		}
		return cards[i].Grade > cards[j].Grade
	})
}

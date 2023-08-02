package util

import (
	"fmt"
	"testing"
	"wan_go/internal/landlord/model"
	"wan_go/pkg/common/constant/landlord_const"
)

func TestEq(t *testing.T) {

	//for i := 1; i < 18; i++ {
	//	if i != 10 && i != 15 {
	//		println(i, false)
	//	} else {
	//		println(i, true)
	//	}
	//}

	for i := 0; i < 18; i++ {
		if i != 8 && i != 12 && i != 16 {
			println(i, false)
		} else {
			println(i, true)
		}
	}

}

func TestIsAircraftWithSingleWing(t *testing.T) {
	// 333444 56	    3,4     5,6
	// 333444 55	    3,4		5,5
	// 333444555 667  3,4,5   6,6,7
	// 333444555 999  3,4,5  ,9
	// 33334444
	// 333444555666 7777

	var cards []*model.Card
	arr0 := []int{3, 3, 4, 4, 4, 5, 6, 3, 2, 2, 2, 1}
	arr1 := []int{3, 3, 4, 4, 5, 4, 5, 3}
	arr2 := []int{3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 7}
	arr3 := []int{3, 3, 3, 4, 4, 4, 5, 5, 5, 9, 9, 9}
	arr4 := []int{3, 3, 3, 3, 4, 4, 4, 4}
	arr5 := []int{3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7, 7}

	for _, i := range arr0 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ := IsAircraftWithSingleWing(cards)
	fmt.Println("arr0:", b)

	cards = make([]*model.Card, 0)
	for _, i := range arr1 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ = IsAircraftWithSingleWing(cards)
	fmt.Println("arr1:", b)

	cards = make([]*model.Card, 0)
	for _, i := range arr2 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ = IsAircraftWithSingleWing(cards)
	fmt.Println("arr2:", b)

	cards = make([]*model.Card, 0)
	for _, i := range arr3 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ = IsAircraftWithSingleWing(cards)
	fmt.Println("arr3:", b)

	cards = make([]*model.Card, 0)
	for _, i := range arr4 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ = IsAircraftWithSingleWing(cards)
	fmt.Println("arr4:", b)

	cards = make([]*model.Card, 0)
	for _, i := range arr5 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ = IsAircraftWithSingleWing(cards)
	fmt.Println("arr5:", b)

}

func TestIsAircraftWithPairWing(t *testing.T) {
	// 333444555667788
	// 333444555666677
	// 3334445566
	// 3334445555

	var cards []*model.Card
	arr0 := []int{3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 7, 7, 8, 8}
	arr1 := []int{3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 8, 8}
	arr2 := []int{3, 3, 3, 4, 4, 4, 5, 5, 6, 6}
	arr3 := []int{3, 3, 3, 4, 4, 4, 5, 5, 5, 5}
	arr4 := []int{3, 3, 3, 3, 4, 4, 4, 4, 5, 5}
	arr5 := []int{3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7}

	for _, i := range arr0 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ := IsAircraftWithPairWing(cards)
	fmt.Println("arr0:", b)

	cards = make([]*model.Card, 0)
	for _, i := range arr1 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ = IsAircraftWithPairWing(cards)
	fmt.Println("arr1:", b)

	cards = make([]*model.Card, 0)
	for _, i := range arr2 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ = IsAircraftWithPairWing(cards)
	fmt.Println("arr2:", b)

	cards = make([]*model.Card, 0)
	for _, i := range arr3 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ = IsAircraftWithPairWing(cards)
	fmt.Println("arr3:", b)

	cards = make([]*model.Card, 0)
	for _, i := range arr4 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	b, _ = IsAircraftWithPairWing(cards)
	fmt.Println("arr4:", b)

	cards = make([]*model.Card, 0)
	for _, i := range arr5 {
		cards = append(cards, &model.Card{Grade: landlord_const.CardGrade(i)})
	}
	fmt.Println("arr5:", b)

}

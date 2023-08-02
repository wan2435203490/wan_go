package util

import (
	"fmt"
	"testing"
	"wan_go/internal/landlord/model"
)

func Test_sortCards(t *testing.T) {
	cards := []*model.Card{
		{Grade: 3},
		{Grade: 2},
		{Grade: 6},
		{Grade: 8},
		{Grade: 1},
		{Grade: 2},
		{Grade: 3},
		{Grade: 5},
		{Grade: 3},
		{Grade: 7},
	}

	sortCards(cards)

	for _, card := range cards {
		fmt.Println(card.Grade)
	}

}

func getCards() []*model.Card {
	return []*model.Card{
		{Grade: 3},
		{Grade: 10},
		{Grade: 11},
		{Grade: 2},
		{Grade: 2},
		{Grade: 6},
		{Grade: 5},
		{Grade: 8},
		{Grade: 8},
		{Grade: 8},
		{Grade: 1},
		{Grade: 3},
		{Grade: 6},
		{Grade: 6},
		{Grade: 5},
		{Grade: 7},
		{Grade: 7},
		{Grade: 7},
		{Grade: 9},
		{Grade: 5},
		{Grade: 5},
		{Grade: 13},
		{Grade: 12},
		{Grade: 14},
	}

}

func TestGetSingleCards(t *testing.T) {
	cards := GiveSingleCards(getCards(), 5)
	printCards(cards)
}

func TestGetPairCards(t *testing.T) {
	cards := GivePairCards(getCards(), 3)
	printCards(cards)
}

func TestGetThreeCards(t *testing.T) {
	cards := GiveThreeCards(getCards(), 3)
	printCards(cards)
}

func TestGiveFourCards(t *testing.T) {
	cards := GiveFourCards(getCards(), 3)
	printCards(cards)
}

func TestGiveStraightCards(t *testing.T) {
	cards := GiveStraightCards(getCards(), 1, 5)
	printCards(cards)
}

func TestGiveStraightPairCards(t *testing.T) {
	cards := GiveStraightPairCards(getCards(), 1, 6)
	printCards(cards)
}

func TestGiveAircraft(t *testing.T) {
	cards := GiveAircraft(getCards(), 1, 9)
	printCards(cards)
}

func printCards(cards [][]*model.Card) {
	for i, card := range cards {
		fmt.Printf("\n第%d次推荐：\n\t", i+1)
		for _, p := range card {
			fmt.Printf("\t%d", p.Grade)
		}
	}
}

func TestAllBomb(t *testing.T) {

	cards := getCards()
	sortCards(cards)
	bombs := findAllBombs(cards)
	fmt.Println(len(bombs))
	for i, card := range bombs {
		fmt.Printf("\n第%d次推荐：\n\t", i+1)
		for _, p := range card {
			fmt.Printf("\t%d", p.Grade)
		}
	}

}

func Test_findBomb(t *testing.T) {

	bombs := GiveBombs(getCards(), -1)
	if bombs == nil {
		fmt.Print("aaaaaa\n")
	}
	fmt.Printf("%#v\n", bombs)
	for _, card := range bombs {
		for _, p := range card {
			fmt.Printf("%#v\n", p)
		}
	}
}

func Test_findJokerBomb(t *testing.T) {
	cards := getCards()
	bomb := findJokerBomb(cards)
	for _, card := range bomb {
		fmt.Printf("%#v\n", card)
	}
}

func TestArr(t *testing.T) {
	var p *int
	p2 := (*int)(nil)
	fmt.Println(p == p2)
}

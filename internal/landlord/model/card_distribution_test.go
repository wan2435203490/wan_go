package model

import (
	"fmt"
	"testing"
)

func TestCreateCards(t *testing.T) {

	distribution := &CardDistribution{}

	distribution.CreateCards()

	fmt.Println("i", "type", "number", "grade")
	for i, card := range distribution.AllCardList {
		fmt.Println(i, card.Type.GetCardType(), card.Number.GetCardNumber(), card.Grade)
	}
}

func TestShuffle(t *testing.T) {
	distribution := &CardDistribution{}

	distribution.CreateCards()

	distribution.Shuffle()

	fmt.Println("i", "type", "number", "grade")
	for i, card := range distribution.AllCardList {
		fmt.Println(i, card.Type.GetCardType(), card.Number.GetCardNumber(), card.Grade)
	}
}

func TestDeal(t *testing.T) {
	distribution := &CardDistribution{}

	distribution.Refresh()

	fmt.Println("i", "type", "number", "grade")
	for i, card := range distribution.AllCardList {
		fmt.Println(i, card.Type.GetCardType(), card.Number.GetCardNumber(), card.Grade)
	}
	fmt.Println("player1---------------------")
	for i, card := range distribution.Player1Cards {
		fmt.Println(i, card.Type.GetCardType(), card.Number.GetCardNumber(), card.Grade)
	}
	fmt.Println("player2---------------------")
	for i, card := range distribution.Player2Cards {
		fmt.Println(i, card.Type.GetCardType(), card.Number.GetCardNumber(), card.Grade)
	}
	fmt.Println("player3---------------------")
	for i, card := range distribution.Player3Cards {
		fmt.Println(i, card.Type.GetCardType(), card.Number.GetCardNumber(), card.Grade)
	}
	fmt.Println("top---------------------")
	for i, card := range distribution.TopCards {
		fmt.Println(i, card.Type.GetCardType(), card.Number.GetCardNumber(), card.Grade)
	}
}

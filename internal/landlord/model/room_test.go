package model

import (
	"fmt"
	"testing"
	"wan_go/pkg/common/constant/landlord_const"
)

func TestReset(t *testing.T) {
	pl := make([]*Player, 0, 8)
	for i := 0; i < 8; i++ {
		pl = append(pl, &Player{
			Identity: landlord_const.Identity(i),
		})
	}

	r := &Room{
		PlayerList: pl,
	}

	r.Reset()

	fmt.Printf("%#v", r)

}

func TestFor(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}

	for i, i2 := range s {
		fmt.Println(i, &i2)
	}
}

func TestPlayerId(t *testing.T) {
	room := &Room{}
	room.PlayerList = []*Player{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}

	id := room.GetAvailablePlayerId()

	fmt.Println(id)
}

package landlord_const

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJson(t *testing.T) {
	var s []CardType
	s = []CardType{Spade, Diamond, Heart, Club, SmallJokerType, BigJokerType, Spade}

	bs, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))

	var out []CardType
	err = json.Unmarshal(bs, &out)
	if err != nil {
		panic(err)
	}
}

func TestFor(t *testing.T) {
	var ids []*int
	for i := 0; i < 10; i++ {
		ids = append(ids, &i)
		//if i == 5 {
		//	i = 8
		//}
	}

	for _, v := range ids {
		fmt.Println(v, *v)
	}
	fmt.Printf("%#v", ids)
}

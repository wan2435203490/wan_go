package model

import (
	"wan_go/pkg/common/constant/landlord_const"
)

func GetCard(id int) *Card {

	var t landlord_const.CardType
	var g landlord_const.CardGrade
	var n landlord_const.CardNumber

	switch id {
	case 52:
		t = landlord_const.SmallJokerType
		g = landlord_const.Fourteenth
		n = landlord_const.SmallJoker
	case 53:
		t = landlord_const.BigJokerType
		g = landlord_const.Fifteenth
		n = landlord_const.BigJoker
	default:
		l := id / 4
		t = landlord_const.CardType(id % 4)
		n = landlord_const.CardNumber(l + 1)
		g = landlord_const.CardGrade((l + 11) % 13)
	}

	card := &Card{
		Type:   t,
		Number: n,
		Grade:  g,
	}

	return card
}

//func GetCard(id int) *Card {
//	t := landlord_const.CardType(id % 4)
//	n := id / 4
//	var g landlord_const.CardGrade
//
//	if n > 12 {
//		if t == 0 {
//			t = landlord_const.SmallJokerType
//			g = landlord_const.Fourteenth
//		} else {
//			t = landlord_const.BigJokerType
//			g = landlord_const.Fifteenth
//		}
//	}
//
//	card := &Card{
//		Type:   t,
//		Number: landlord_const.CardNumber(n + 1),
//	}
//
//	if g != 0 {
//		card.Grade = g
//	} else {
//		card.Grade = ConvertNum2Grade(n)
//	}
//
//	return card
//}
//
//func ConvertNum2Grade(n int) landlord_const.CardGrade {
//	//n: 0,1,2...13
//	//n:0, Grade:A Twelfth(11)
//	//n:1, Grade:2 Thirteenth(12)
//	//n:2, Grade:3 First(0)
//	//(过滤掉n:13了)n:13, Grade:小王 Fourteenth(14) ｜ 大王 Fifteenth(15)
//	if n == 13 {
//		panic("你咋进来的呢")
//	}
//
//	return landlord_const.CardGrade((n + 11) % 13)
//}

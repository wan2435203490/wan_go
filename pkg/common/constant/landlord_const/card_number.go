package landlord_const

type CardNumber int

const (
	DiscardCardNumber CardNumber = iota
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten

	Jack
	Lady
	King

	SmallJoker
	BigJoker
)

func (cn CardNumber) GetCardNumber() int {
	return []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}[cn]
}

//
//func (c *CardNumber) UnmarshalJSON(b []byte) error {
//	var s int
//	if err := json.Unmarshal(b, &s); err != nil {
//		return err
//	}
//	switch strings.ToUpper(s) {
//	default:
//		*c = One
//	case "ONE":
//		*c = One
//	case "TWO":
//		*c = Two
//	case "THREE":
//		*c = Three
//	case "FOUR":
//		*c = Four
//	case "FIVE":
//		*c = Five
//	case "SIX":
//		*c = Six
//	case "SEVEN":
//		*c = Seven
//	case "EIGHT":
//		*c = Eight
//	case "NINE":
//		*c = Nine
//	case "TEN":
//		*c = Ten
//	case "JACK":
//		*c = Jack
//	case "LADY":
//		*c = Lady
//	case "KING":
//		*c = King
//	case "SMALL_JOKER":
//		*c = SmallJoker
//	case "BIG_JOKER":
//		*c = BigJoker
//	}
//
//	return nil
//}
//
//func (c CardNumber) MarshalJSON() ([]byte, error) {
//	var s string
//	switch c {
//	default:
//		s = "ONE"
//	case One:
//		s = "ONE"
//	case Two:
//		s = "TWO"
//	case Three:
//		s = "THREE"
//	case Four:
//		s = "FOUR"
//	case Five:
//		s = "FIVE"
//	case Six:
//		s = "SIX"
//	case Seven:
//		s = "SEVEN"
//	case Eight:
//		s = "EIGHT"
//	case Nine:
//		s = "NINE"
//	case Ten:
//		s = "TEN"
//	case Jack:
//		s = "JACK"
//	case Lady:
//		s = "LADY"
//	case King:
//		s = "KING"
//	case SmallJoker:
//		s = "SMALL_JOKER"
//	case BigJoker:
//		s = "BIG_JOKER"
//	}
//
//	return json.Marshal(s)
//}

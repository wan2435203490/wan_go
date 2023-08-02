package landlord_const

import (
	"encoding/json"
	"strings"
)

type Type int

const (
	Single Type = iota
	Pair
	ThreeType
	ThreeWithOne
	ThreeWithPair

	FourWithTwo
	FourWithFour

	Bomb
	JokerBomb

	Straight
	StraightPair

	Aircraft
	AircraftWithSingleWings
	AircraftWithPairWings
)

func (t Type) GetType() string {
	return []string{"单张", "对子", "三张", "三带一", "三带一对", "四带二", "四带两对", "炸弹", "王炸", "顺子",
		"连对", "飞机", "飞机带翅膀", "飞机带大炮"}[t]
}

func (t *Type) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToUpper(s) {
	default:
		*t = Single
	case "SINGLE":
		*t = Single
	case "PAIR":
		*t = Pair
	case "THREE":
		*t = ThreeType
	case "THREE_WITH_ONE":
		*t = ThreeWithOne
	case "THREE_WITH_PAIR":
		*t = ThreeWithPair
	case "FOUR_WITH_TWO":
		*t = FourWithTwo
	case "FOUR_WITH_FOUR":
		*t = FourWithFour
	case "BOMB":
		*t = Bomb
	case "JOKER_BOMB":
		*t = JokerBomb
	case "STRAIGHT":
		*t = Straight
	case "STRAIGHT_PAIR":
		*t = StraightPair
	case "AIRCRAFT":
		*t = Aircraft
	case "AIRCRAFT_WITH_SINGLE_WINGS":
		*t = AircraftWithSingleWings
	case "AIRCRAFT_WITH_PAIR_WINGS":
		*t = AircraftWithPairWings
	}

	return nil
}

func (t Type) MarshalJSON() ([]byte, error) {
	var s string
	switch t {
	default:
		s = "SINGLE"
	case Single:
		s = "SINGLE"
	case Pair:
		s = "PAIR"
	case ThreeType:
		s = "THREE"
	case ThreeWithOne:
		s = "THREE_WITH_ONE"
	case ThreeWithPair:
		s = "THREE_WITH_PAIR"
	case FourWithTwo:
		s = "FOUR_WITH_TWO"
	case FourWithFour:
		s = "FOUR_WITH_FOUR"
	case Bomb:
		s = "BOMB"
	case JokerBomb:
		s = "JOKER_BOMB"
	case Straight:
		s = "STRAIGHT"
	case StraightPair:
		s = "STRAIGHT_PAIR"
	case Aircraft:
		s = "AIRCRAFT"
	case AircraftWithSingleWings:
		s = "AIRCRAFT_WITH_SINGLE_WINGS"
	case AircraftWithPairWings:
		s = "AIRCRAFT_WITH_PAIR_WINGS"
	}

	return json.Marshal(s)
}

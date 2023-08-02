package landlord_const

type DimensionType int

const (
	All DimensionType = iota
	Room
)

func (d DimensionType) DimensionType() string {
	return [2]string{"ALL", "ROOM"}[d]
}

func ToDimensionType(str string) DimensionType {

	switch str {
	case "ALL":
		return All
	case "ROOM":
		return Room
	default:
		return Room
	}
}

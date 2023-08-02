package landlord_const

import (
	"encoding/json"
	"strings"
)

type Identity int

const (
	Default Identity = iota
	Landlord
	Farmer
)

func (i Identity) GetIdentity() string {
	return []string{"", "地主", "农民"}[i]
}

func (i *Identity) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToUpper(s) {
	case "LANDLORD":
		*i = Landlord
	case "FARMER":
		*i = Farmer
	default:
		*i = Default
	}

	return nil
}

func (i Identity) MarshalJSON() ([]byte, error) {
	var s string
	switch i {
	case Landlord:
		s = "LANDLORD"
	case Farmer:
		s = "FARMER"
	default:
		s = ""
	}

	return json.Marshal(s)
}

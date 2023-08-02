package model

import "wan_go/pkg/common/constant/landlord_const"

type RoundResult struct {
	WinIdentity landlord_const.Identity `json:"winIdentity"`
	LandlordId  int                     `json:"landlord"`
	Multiple    int                     `json:"multiple"`
}

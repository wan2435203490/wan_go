package dto

import (
	"strconv"
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/utils"
)

type ResultScore struct {
	UserName     string `json:"username"`
	MoneyChange  string `json:"moneyChange"`
	IdentityName string `json:"identityName"`
}

func NewResultScore(userName string, multiple int, isWin, isLandlord bool) *ResultScore {
	if isLandlord {
		return NewLandlordScore(userName, multiple, isWin)
	} else {
		return NewFarmerScore(userName, multiple, isWin)
	}
}

func NewFarmerScore(userName string, multiple int, isWin bool) *ResultScore {
	return &ResultScore{
		UserName:     userName,
		MoneyChange:  GetMoneyChange(multiple, isWin),
		IdentityName: landlord_const.Farmer.GetIdentity(),
	}
}

func NewLandlordScore(userName string, multiple int, isWin bool) *ResultScore {
	return &ResultScore{
		UserName:     userName,
		MoneyChange:  GetMoneyChange(multiple*2, isWin),
		IdentityName: landlord_const.Landlord.GetIdentity(),
	}
}

func GetMoneyChange(moneyChange int, isWin bool) string {
	moneyChange = utils.IfThen(isWin, moneyChange, -moneyChange).(int)
	return strconv.Itoa(moneyChange)
}

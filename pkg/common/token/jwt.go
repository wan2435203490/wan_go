package token

import (
	"github.com/golang-jwt/jwt/v4"
	"wan_go/common/config"
	"wan_go/common/constant"
	"wan_go/common/utils"
	"time"
)

type Claims struct {
	UserId string
	jwt.RegisteredClaims
}

func BuildClaims(userId string, day int64) Claims {
	now := time.Now()
	before := now.Add(-time.Minute * 5)
	return Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  []string{"cy"},
			Issuer:    "cy",
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(day*24) * time.Hour)), //Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                                        //Issuing time
			NotBefore: jwt.NewNumericDate(before),                                     //Begin Effective time
		}}
}

func CreateToken(userId string) (string, error) {
	claims := BuildClaims(userId, config.Config.TokenPolicy.JwtExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.TokenPolicy.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.TokenPolicy.JwtSecret), nil
	}
}

func GetClaimFromToken(tokensString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, utils.Wrap(constant.ErrTokenMalformed, "")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, utils.Wrap(constant.ErrTokenExpired, "")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, utils.Wrap(constant.ErrTokenNotValidYet, "")
			} else {
				return nil, utils.Wrap(constant.ErrTokenUnknown, "")
			}
		} else {
			return nil, utils.Wrap(constant.ErrTokenNotValidYet, "")
		}
	} else {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, utils.Wrap(constant.ErrTokenNotValidYet, "")
	}
}

func GetUserIdFromToken(tokenStr string) (bool, string) {
	claims, err := GetClaimFromToken(tokenStr)
	if err != nil {
		errMsg := "parse token err " + err.Error()
		return false, errMsg
	}
	return true, claims.UserId
}

func WsVerifyToken(tokenStr, userId string) (bool, error, string) {
	claims, err := GetClaimFromToken(tokenStr)
	if err != nil {
		errMsg := "parse token err " + err.Error()
		return false, utils.Wrap(err, errMsg), errMsg
	}
	if claims.UserId != userId {
		errMsg := " userId is not same to token userId. claims.UserId: " + claims.UserId
		return false, utils.Wrap(constant.ErrTokenDifferentUserID, errMsg), errMsg
	}

	return true, nil, ""
}

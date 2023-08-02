package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/utils"
	"wan_go/sdk/pkg/jwtauth"
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
	claims := BuildClaims(userId, config.Config.Jwt.Expire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.Jwt.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Jwt.Secret), nil
	}
}

func GetClaimFromToken(tokensString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, utils.Wrap(constant.ErrTokenMalformed)
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, utils.Wrap(constant.ErrTokenExpired)
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, utils.Wrap(constant.ErrTokenNotValidYet)
			} else {
				return nil, utils.Wrap(constant.ErrTokenUnknown)
			}
		} else {
			return nil, utils.Wrap(constant.ErrTokenNotValidYet)
		}
	} else {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, utils.Wrap(constant.ErrTokenNotValidYet)
	}
}

//func GetUserIdFromToken(tokenStr string) (bool, string) {
//	claims, err := GetClaimFromToken(tokenStr)
//	if err != nil {
//		errMsg := "parse token err " + err.Error()
//		return false, errMsg
//	}
//	return true, claims.UserId
//}

func GetJWTToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		//todo config
		if jwt.GetSigningMethod("HS256") != t.Method {
			return nil, jwtauth.ErrInvalidSigningAlgorithm
		}
		//todo
		//if usingPublicKeyAlgo {
		//	return pubKey, nil
		//}
		//c.Set("constant.JWTToken", tokenStr)

		return []byte(config.Config.Jwt.Secret), nil
	})
}

func GetUserIdFromToken(tokenStr string) (int32, error) {
	token, err := GetJWTToken(tokenStr)
	if err != nil {
		return -1, err
	}
	claims := jwtauth.MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		if key == constant.ClaimsIdentity {
			claims[key] = value
			break
		}
	}
	if claims[constant.ClaimsIdentity] != nil {
		return int32((claims[constant.ClaimsIdentity]).(float64)), nil
	}
	return -1, errors.New("identity not exists")
}

func WsVerifyToken(tokenStr, userId string) (bool, error, string) {
	claims, err := GetClaimFromToken(tokenStr)
	if err != nil {
		errMsg := "parse token err " + err.Error()
		return false, utils.WrapMsg(err, errMsg), errMsg
	}
	if claims.UserId != userId {
		errMsg := " userId is not same to token userId. claims.UserId: " + claims.UserId
		return false, utils.WrapMsg(constant.ErrTokenDifferentUserID, errMsg), errMsg
	}

	return true, nil, ""
}

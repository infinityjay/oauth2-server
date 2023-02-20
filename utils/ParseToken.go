package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"oauth2/model"
)

// ParseToken 解析JWT
func ParseToken(tokenString string) (*model.MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString,&model.MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return model.MySecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

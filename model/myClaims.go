package model

import (
	"context"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
	"strings"
)

type MyClaims struct {
	UserName string   `json:"user_name"`
	UserType []string `json:"user_type"`
	jwt.StandardClaims
}

var MySecret = []byte("3023273")

func (mt *MyClaims) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	mt.Audience = data.Client.GetID()
	mt.Subject = data.UserID
	mt.ExpiresAt = data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mt)
	access, err := token.SignedString(MySecret)
	if err != nil {
		return "", "", err
	}
	refresh := ""

	if isGenRefresh {
		t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
		refresh = base64.URLEncoding.EncodeToString([]byte(t))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}
	return access, refresh, err
}

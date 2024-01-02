package goutils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(secret string, data string) (string, error) {
	//time.Now().Add(time.Minute * 15)
	var err error
	//Creating Access Token

	atClaims := jwt.MapClaims{}
	atClaims["data"] = data
	//atClaims["user_id"] = userid
	//atClaims["exp"] = expiresin.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(secret string, token string) (string, error) {

	//map[string]interface{}
	atClaims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &atClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Printf("%#v\n", err)
		return "", err
	}
	return atClaims["data"].(string), nil
}

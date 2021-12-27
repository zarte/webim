package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"websocket/Config"
)

/**
解析token
*/
func ParseToken(tokenString string)(*CustomClaims,error)  {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(Config.Gconfig.SecretKey), nil
	})
	if token==nil{
		return nil,fmt.Errorf("token 异常")
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims,nil
	} else {
		return nil,err
	}
}

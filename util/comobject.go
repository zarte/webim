package util

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	UserId string
	jwt.StandardClaims
}
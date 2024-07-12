package types

import "github.com/golang-jwt/jwt/v5"

type UserJwtClaims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

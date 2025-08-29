package dto

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	Role       string   `json:"role"`
	Permission []string `json:"permissions"`
	jwt.RegisteredClaims
}

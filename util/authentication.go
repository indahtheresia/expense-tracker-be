package util

import (
	"errors"
	"expense-tracker/constant"
	"expense-tracker/dto"

	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateBcrypt(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, constant.ErrorGenerateHash
	}

	passHash := string(hash)

	return &passHash, nil
}

func CompareHashPassword(hashPassword []byte, userPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashPassword, userPassword)
	if err != nil {
		return constant.ErrorComparePassword
	}
	return nil
}

func GenerateJWTToken(regClaims *dto.CustomClaims) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, regClaims)

	key := os.Getenv("SIGNEDPASS")

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

var ParseToken = func(tokenString string) (*dto.CustomClaims, error) {
	var myCustomClaims dto.CustomClaims

	token, err := jwt.ParseWithClaims(tokenString, &myCustomClaims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("SIGNEDPASS")), nil
	},
		jwt.WithIssuedAt(),
		jwt.WithIssuer(os.Getenv("APP_NAME")),
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, constant.ErrorHandleTokenExpired
		}
		return nil, fmt.Errorf("error: %v", err)
	}

	claims, ok := token.Claims.(*dto.CustomClaims)
	if !ok || !token.Valid {
		return nil, constant.ErrorUnknownClaimsType
	}

	return claims, nil
}

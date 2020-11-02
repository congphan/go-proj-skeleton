package auth

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

type (
	Principal struct {
		ID string
	}

	accessClaims struct {
		UserID string `json:"user_id"`
		jwt.StandardClaims
	}

	TokenValidator struct {
		secret []byte
	}
)

func NewTokenValidator(secretKey []byte) *TokenValidator {
	return &TokenValidator{
		secretKey,
	}
}

func (v TokenValidator) Validate(accessToken string) (*Principal, error) {
	claims := &accessClaims{}
	parsedToken, err := jwt.ParseWithClaims(accessToken, claims, v.keyVerifier)
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.UserID == "" {
		return nil, fmt.Errorf("missing claims [user_id]")
	}

	return &Principal{
		claims.UserID,
	}, nil
}

func (v TokenValidator) keyVerifier(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
	}

	return v.secret, nil
}

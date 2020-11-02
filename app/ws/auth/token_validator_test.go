package auth

import (
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func Test_TokenValidator_Validate(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		t.Run("use different key", func(t *testing.T) {
			secret := []byte("secret_key")

			claims := accessClaims{
				UserID: "123456",
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signedToken, err := token.SignedString(secret)
			assert.NoError(t, err)

			validator := TokenValidator{
				[]byte("secret"),
			}
			_, err = validator.Validate(signedToken)
			assert.EqualError(t, err, "signature is invalid")
		})

		t.Run("unexpected signing method", func(t *testing.T) {
			secret := []byte("secret_key")

			claims := accessClaims{
				UserID: "123456",
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
			signedToken, err := token.SignedString(secret)
			assert.NoError(t, err)

			validator := TokenValidator{
				[]byte("secret_key"),
			}
			_, err = validator.Validate(signedToken)
			assert.EqualError(t, err, "Unexpected signing method: HS384")
		})

		t.Run("time expires", func(t *testing.T) {
			secret := []byte("secret_key")

			claims := accessClaims{
				UserID: "123456",
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: 1592997120, //06/24/2020 @ 11:12am (UTC)
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signedToken, err := token.SignedString(secret)
			assert.NoError(t, err)

			validator := TokenValidator{
				secret,
			}
			_, err = validator.Validate(signedToken)
			if assert.Error(t, err) {
				assert.True(t, strings.Contains(err.Error(), "token is expired"))
			}
		})

		t.Run("missing user_id", func(t *testing.T) {
			secret := []byte("secret_key")

			claims := accessClaims{
				StandardClaims: jwt.StandardClaims{
					Id: "1234",
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signedToken, err := token.SignedString(secret)
			assert.NoError(t, err)

			validator := TokenValidator{
				secret,
			}
			_, err = validator.Validate(signedToken)
			assert.EqualError(t, err, "missing claims [user_id]")
		})
	})

	t.Run("valid", func(t *testing.T) {
		secret := []byte("secret_key")

		claims := accessClaims{
			UserID: "123456",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(secret)
		assert.NoError(t, err)

		expectedSignedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDU2In0.JucMlcxvcClLoFKvZaLygMvDqueUkgaW-SZ9xlrBZgo"
		assert.Equal(t, expectedSignedToken, signedToken)

		validator := TokenValidator{
			secret,
		}
		p, err := validator.Validate(signedToken)
		assert.NoError(t, err)
		assert.Equal(t, "123456", p.ID)
	})
}

package lib

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var TokenSecretKey = []byte("SecretYouShouldHide")

type UserToken struct {
	Username string `json:"usernmae"`
	Email    string `json:"email"`
}

type Claims struct {
	UserToken
	jwt.RegisteredClaims
}

func CreateNewToken(userToken UserToken) (string, error) {
	expirationTime := time.Now().Add(45 * time.Minute)
	claims := &Claims{
		UserToken: userToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(TokenSecretKey)

	if err != nil {
		log.Fatal(err)
		return "", nil
	}

	return token, nil

}

func ValidateToken(token string) bool {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return TokenSecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return false
		}
	}

	return tkn.Valid
}

// Warning it doesn't validate the token expiration
func GetUnverifiedClaims(token string) (Claims, error) {
	parser := jwt.NewParser()
	jwtToken, _, err := parser.ParseUnverified(token, &Claims{})
	if err != nil {
		return Claims{}, err
	}
	claims := jwtToken.Claims.(*Claims)

	return *claims, nil
}

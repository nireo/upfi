package lib

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("something_very_secret")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(username string) (string, error) {
	expirationTime := time.Now().Add(time.Date)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Takes in a token strings, and returns the username in the token.
// also checks for different errors
func ValidateToken(tokenString string) (string, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("Unauthorized")
		}
		return "", errors.New("Bad request")
	}
	if !tkn.Valid {
		return "", errors.New("Token is invalid")
	}

	return claims.Username
}

// Function that gets called during the initial startup so we can easily load ENV variables once
func SetJWTKey(key string) {
	jwtKey = []byte(key)
}

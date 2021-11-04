package lib

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("something_very_secret")

// C is a simple struct to keep the username and the standard claims of a jwt, such as the expiration time
type C struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// CreateToken creates a jwt token which stores a username and is valid for 24 hours.
func CreateToken(username string) (string, error) {
	// Set the expiration time of the token to be 24 hours.
	expirationTime := time.Now().Add(time.Hour * 24)

	// Construct the jsonwebtoken claims
	claims := &C{
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

// ValidateToken takes a token as an argument and checks if that token is valid.
// If the token is valid, then the function returns the usernanem stored in the token.
func ValidateToken(tokenString string) (string, error) {
	claims := &C{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// Check for different errors with the token.
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("Unauthorized")
		}
		return "", errors.New("bad request")
	}
	if !tkn.Valid {
		return "", errors.New("token is invalid")
	}

	// Return the username stored in the token.
	return claims.Username, nil
}

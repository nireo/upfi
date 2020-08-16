package middleware

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var secretKey = []byte("thisismysecret")

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return secretKey, nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}


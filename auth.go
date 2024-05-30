package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func withJWTAuth(handlerFunc http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Get the token from request (Auth Header )
		tokenString := GetTokenFromRequest(r)
		//Validate the token
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to authenticate token ")
			permissionDenied(w)

			return
		}

		if !token.Valid {
			log.Printf("failed to authenticate token ")
			permissionDenied(w)

			return

		}
		//get the userID from the token
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)

		_, err = store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user")
			permissionDenied(w)

			return
		}
		//call the handler func and continue to the endpoint
		handlerFunc(w, r)

	}
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, ErrorResponse{
		Error: fmt.Errorf("permission denied").Error(),
	})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}
	return ""
}

func validateJWT(t string) (*jwt.Token, error) {

	secret := Envs.JWTSecret
	return jwt.Parse(t , func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", t.Header["alg]"])

		}
		return []byte(secret), nil

	})
}

func HashPassword(pw string ) (string,error) { 
	
hash , err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
if err !=nil {
	return "",err
}
return string(hash), nil
}

func CreateJWT(secret[]byte, userID int64 )(string,error){
	token := jwt.NewWithClaims(jwt.SigningMethodES256,jwt.MapClaims{
		"userID" : strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(time.Hour * 24).Unix(),

	})
	tokenString , err := token.SignedString(secret)
	if err!= nil {
		return "",err
	}
	return tokenString,nil

}
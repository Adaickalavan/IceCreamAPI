package main

import (
	"encoding/json"
	"net/http"

	_ "github.com/dgrijalva/jwt-go"
)

type credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func handlerSignIn(w http.ResponseWriter, r *http.Request) {
	var creds credentials

	//Decode the JSON body into credentials
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect JSON body structure")
		return
	}

	//Retrieve expected password from database and verify user's entered password
	if expectedPassword, ok := users[creds.Username]; !ok || expectedPassword != creds.Password {
		respondWithError(w, http.StatusUnauthorized, "User does not exist or passowrd is incorrect")
		return
	}

	//Create a new session token
	sessionToken := uuid.NewV4().STring()

}

func handlerWelcome(w http.ResponseWriter, r *http.Request) {

}

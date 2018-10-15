package main

import (
	"credentials"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//Contains a list of authorised "username:credentials.LoginHash"
var adminDatabase = map[string]credentials.LoginHash{
	"user1": credentials.LoginHash{
		Name:      "user1",
		HashedPwd: []byte("$2a$08$67iAYpKAVCygeyf1mqzOzueitZw.Umk/HczdsLm16Qi547/gbbgg."),
	},
}

type customClaims struct {
	admin              bool   //Set administrator rights
	name               string //Set user name
	jwt.StandardClaims        //Standard JWT claims
}

// Verify the login info with a database and obtain the
// claims(i.e., payload) for this particular webpage and user
func verify(login credentials.Login) (jwt.Claims, error) {
	if loginHash, ok := adminDatabase[login.Name]; !ok || !login.Compare(loginHash) {
		return nil, fmt.Errorf("Username and/or password is incorrect")
	}

	//Form the claims (payload) of the token
	//Do not include sensitive info such as password in the claims
	claims := customClaims{
		admin: true,
		name:  login.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(), //Token validity period
			Issuer:    "HomeBase",
		},
	}

	return claims, nil
}

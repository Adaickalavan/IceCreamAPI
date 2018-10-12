package main

import (
	"credentials"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//Contains a list of authorised "username:password"
var adminsDatabase = map[string]string{
	"user1": "1234",
}

type customClaims struct {
	admin              bool   //Set administrator rights
	name               string //Set user name
	jwt.StandardClaims        //Standard JWT claims
}

// Verify the login info with a database and obtain the
// claims(i.e., payload) for this particular webpage and user
func login(loginInfo credentials.LoginInfo) (jwt.Claims, error) {
	//Verify the login info with a database
	if expectedPassword, ok := adminsDatabase[loginInfo.Name]; !ok || (expectedPassword != loginInfo.Password) {
		return nil, fmt.Errorf("Username and/or password is incorrect")
	}

	//Form the claims (payload) of the token
	//Do not include sensitive info such as password in the claims
	claims := customClaims{
		admin: true,
		name:  loginInfo.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(), //Token validity period
			Issuer:    "HomeBase",
		},
	}

	return claims, nil
}

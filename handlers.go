package main

import (
	"credentials"
	"document"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/definition/", handlerGetDocByID).Methods("GET")
	muxRouter.HandleFunc("/definition", handlerGetDoc).Methods("GET")
	muxRouter.HandleFunc("/definition", handlerPostDoc).Methods("POST")
	muxRouter.HandleFunc("/definition", handlerPutDoc).Methods("PUT")
	muxRouter.HandleFunc("/definition", handlerDeleteDoc).Methods("DELETE")
	muxRouter.HandleFunc("/login", createToken)
	return muxRouter
}

//Retrieve all documents from database
func handlerGetDoc(w http.ResponseWriter, r *http.Request) {
	docs, err := product.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, docs)
}

//Retrieve only document matching query
func handlerGetDocByID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	doc, err := product.FindByValue(query.Get("doc"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, doc)
}

//Post document to database
func handlerPostDoc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var doc document.Icecream
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&doc); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	doc.ID = bson.NewObjectId()
	err := product.Insert(doc)
	switch {
	case mgo.IsDup(err):
		respondWithError(w, http.StatusConflict, err.Error())
		return
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, doc)
}

//Update document in database
func handlerPutDoc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var doc document.Icecream
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&doc); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := product.Update(doc); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusAccepted, map[string]string{"Result": "Successfully updated"})
}

//Delete document from database
func handlerDeleteDoc(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	err := product.Delete(query.Get("doc"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusAccepted, map[string]string{"Result": "Successfully deleted"})
}

//HTTP reply with error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

//HTTP reply with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		http.Error(w, "HTTP 500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}

type customClaims struct {
	admin              bool              //Set administrator rights
	user               *credentials.User //Set `User` properties
	jwt.StandardClaims                   //Standard JWT claims
}

//Sign JWT with secret signingKey
var signingKey = "signJwtUsingSecretKey"

//Creates JSON web token for users
func createToken(w http.ResponseWriter, r *http.Request) {
	//Decode client response into User struct
	var user credentials.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	//Form the claims (payload) of the token
	claims := customClaims{
		admin: true,
		user:  &user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 80).Unix(), //Token validity period
			Issuer:    "icecreamapi",
		},
	}

	//Create JWT with signing method and claims(i.e. payload)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, //Type of jwt.SigningMethodHS256 is *jwt.SigningMethodHMAC
		claims)

	fmt.Println("Token struct==", token)
	fmt.Println("Token struct==", &token)

	// Sign token with key
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to sign token")
		return
	}

	//Respond with the token
	respondWithJSON(w, http.StatusCreated, credentials.JWToken{TokenString: tokenString})
}

//Authenticate prodives authentication middleware for handlers
func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		//Get token from authorization header
		//format = map[string][]string
		//entry  = "Authorization": "Bearer <token>"
		tokens, ok := r.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			tokenString = tokens[0]
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		//Token is empty
		if tokenString == "" {
			respondWithError(w, http.StatusUnauthorized, "Empty token")
			return
		}

		//Parse takes the token string and a function for looking up the key.
		parsedToken, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {
				//Verify the signing algorithm by using interface type assertion
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return signingKey, nil
			},
		)
		if err != nil {
			value, ok := err.(*jwt.ValidationError)
			switch {
			case ok && (value.Errors&jwt.ValidationErrorMalformed != 0):
				respondWithError(w, http.StatusUnauthorized, "This is not a token")
			case ok && (value.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0):
				respondWithError(w, http.StatusUnauthorized, "Token is expired or not valid yet")
			default:
				respondWithError(w, http.StatusUnauthorized, "Token parsing error: "+err.Error())
			}
			return
		}

		//Check token validity and token claims, and call the `next` function
		if _, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			next(w, r)
			return
		}
		respondWithError(w, http.StatusUnauthorized, "Token claims are not dechiperable and/or token is invalid")
		return
	})
}

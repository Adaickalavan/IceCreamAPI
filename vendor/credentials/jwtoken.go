package credentials

import (
	"encoding/json"
	"fmt"
	"handler"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// Reference: https://godoc.org/github.com/dgrijalva/jwt-go#Token
// type Token struct {
//     Raw       string                 // The raw token.  Populated when you Parse a token
//     Method    SigningMethod          // The signing method used or to be used
//     Header    map[string]interface{} // The first segment of the token
//     Claims    Claims                 // The second segment of the token
//     Signature string                 // The third segment of the token.  Populated when you Parse a token
//     Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
// }

// JWToken struct contains jwt token
type jwToken struct {
	TokenString string `json:"tokenString"`
}

//Sign JWT with secret signingKey
var signingKey = []byte("signJwtUsingSecretKey")

//CreateToken creates JSON web token for users
func CreateToken(next func(Login) (jwt.Claims, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Decode client response into LoginInfo struct
		var login Login
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&login); err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Verify the login info with a database and obtain the
		// claims(i.e., payload) for this particular webpage and user
		claims, err := next(login)
		if err != nil {
			handler.RespondWithError(w, http.StatusUnauthorized, "Failed to login: "+err.Error())
			return
		}

		//Create JWT with signing method and claims(i.e. payload)
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256, //Type of jwt.SigningMethodHS256 is *jwt.SigningMethodHMAC
			claims)

		// Sign token with key
		tokenString, err := token.SignedString(signingKey)
		if err != nil {
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to sign token")
			return
		}

		//Respond with the token
		handler.RespondWithJSON(w, http.StatusCreated, jwToken{TokenString: tokenString})
	})
}

//Authenticate prodives authentication middleware for handlers
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
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
			handler.RespondWithError(w, http.StatusUnauthorized, "Empty token")
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
				handler.RespondWithError(w, http.StatusUnauthorized, "This is not a token")
			case ok && (value.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0):
				handler.RespondWithError(w, http.StatusUnauthorized, "Token is expired or not valid yet")
			default:
				handler.RespondWithError(w, http.StatusUnauthorized, "Token parsing error: "+err.Error())
			}
			return
		}

		//Check token validity and token claims, and call the `next` function
		if _, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			next(w, r)
			return
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Token claims are not dechiperable and/or token is invalid")
		return
	})
}

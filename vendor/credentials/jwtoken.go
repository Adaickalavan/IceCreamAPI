package credentials

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
type JWToken struct {
	TokenString string `json:"tokenString"`
}


package credentials

// JWToken struct contains jwt token
type JWToken struct {
	Token string `json:"token"`
}

// Exception struct contains message to be thrown at exception
type Exception struct {
	Message string `json:"message"`
}

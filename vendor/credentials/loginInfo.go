package credentials

// LoginInfo struct contains user authentication details
type LoginInfo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

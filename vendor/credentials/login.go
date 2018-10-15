package credentials

import "golang.org/x/crypto/bcrypt"

// Login struct contains user authentication details
type Login struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LoginHash struct contains hashed user authentication details
type LoginHash struct {
	Name      string `json:"name"`
	HashedPwd []byte `json:"hashedPwd"`
}

// Hash the password with salt using the bcrypt algorithm
func (login *Login) Hash() (LoginHash, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(login.Password), 8) //Second argument is the cost of hashing, which is arbitrarily set as 8
	if err != nil {
		return LoginHash{}, err
	}
	return LoginHash{Name: login.Name, HashedPwd: hashedPwd}, nil
}

// Compare compares password with hashed password
func (login *Login) Compare(loginHash LoginHash) bool {
	// Salt and hash the password using the bcrypt algorithm
	// Second argument is the cost of hashing, which is arbitrarily set as 8
	err := bcrypt.CompareHashAndPassword(loginHash.HashedPwd, []byte(login.Password))
	if err != nil {
		return false
	}
	return true
}

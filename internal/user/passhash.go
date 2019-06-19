package user

import (
	"golang.org/x/crypto/bcrypt"
)

// https://gowebexamples.com/password-hashing/
// https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

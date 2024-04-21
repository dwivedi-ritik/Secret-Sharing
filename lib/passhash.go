package lib

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CreatePasswordHash(plainText string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(hashed)
}

func ComparePasswordHash(hashed string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

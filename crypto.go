package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	cost := 14

	salt, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return ""
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(salt, cost)
	if err != nil {
		return ""
	}

	return string(salt) + string(hashedPassword)
}

func CheckPasswordHash(password, hash string) bool {
	hashBytes := []byte(hash)

	err := bcrypt.CompareHashAndPassword(hashBytes, []byte(password))
	fmt.Println(err)

	return err == nil
}

package auth

import (
	"log"
	"github.com/alexedwards/argon2id"
)


//Create a hash to store on bd
func HashPassword(password string) (string, error) {

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return hash, nil
}

//Check the password stored on bd with the password login
func CheckPasswordHash(password, hash string) (bool, error) {

	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	log.Printf("Match: %v", match)
	if match == false {
		return false, nil
	}
	return true, nil
}

package bcrypt_util

import "golang.org/x/crypto/bcrypt"

func Compare(input_string string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input_string))
	return err == nil
}

func Hash(input_string string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input_string), bcrypt.DefaultCost)
	return string(hash), err
}

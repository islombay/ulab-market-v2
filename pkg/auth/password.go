package auth_lib

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

func GetHashPassword(s string) (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwd), nil
}

func GenerateRandomPassword(codeLength int) string {
	rand.Seed(time.Now().UnixNano())

	var result string
	for i := 0; i < codeLength; i++ {
		result += strconv.Itoa(rand.Intn(10))
	}
	return result
}

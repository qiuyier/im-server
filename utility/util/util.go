package util

import "fmt"

func ValidatePassword(username, password string) bool {
	encryptPassword := fmt.Sprintf("%s_test", username)
	if password == encryptPassword {
		return true
	} else {
		return false
	}
}

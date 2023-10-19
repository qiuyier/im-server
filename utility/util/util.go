package util

import (
	"fmt"
	"unicode"
)

func ValidatePassword(username, password string) bool {
	encryptPassword := fmt.Sprintf("%s_test", username)
	if password == encryptPassword {
		return true
	} else {
		return false
	}
}

func IsBlank(str string) bool {
	strLen := len(str)
	if str == "" || strLen == 0 {
		return true
	}
	for i := 0; i < strLen; i++ {
		if unicode.IsSpace(rune(str[i])) == false {
			return false
		}
	}
	return true
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

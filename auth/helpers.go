package auth

import (
	"fmt"
	"regexp"
)

const MinimumPasswordLength = 5

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// isEmailValid checks if the email provided passes the required structure and length.
func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func IsIDValid(id int) bool {
	fmt.Println(id)
	if id <= 0 {
		fmt.Println("id less than 0 ", id)
		return false
	}
	return true
}

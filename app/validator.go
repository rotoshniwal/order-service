package app

import (
	"strconv"
)

//isNumber validates if an input string is a valid number or not.
//Returns true if it is a valid number, else false.
func isNumber(id string) bool {
	if _, err := strconv.Atoi(id); err == nil {
		return true
	} else {
		return false
	}
}

//isEmpty validates if a string is empty or not.
func isEmpty(str string) bool {
	if str == "" || len(str) == 0 {
		return true
	} else {
		return false
	}
}

//isEAN validates if an input string (Product code) is in valid EAN-13 format or not.
//Returns true if string is EAN-13 compliant, else false.
func isEAN(str string) bool {
	if isNumber(str) && len(str) == 13 {
		return true
	} else {
		return false
	}
}

package services

import (
	"regexp"
	"unicode"
)

func IsPhoneValid(phone string) bool {
	phoneRegex := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	return phoneRegex.MatchString(phone)
}

func IsPasswordValid(pwd string) bool {
	symbols := 0
	number := false
	upper := false
	special := false
	for _, c := range pwd {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
		default:
			//return
		}
		symbols++
	}

	if symbols >= 8 && number && upper && special {
		return true
	} else {
		return false
	}
}

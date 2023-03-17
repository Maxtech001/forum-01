package registration

import (
	"unicode"
	"regexp"
)

type Password struct {
	Lowercase bool
	Uppercase bool
	Number    bool
	Special   bool
	NoSpaces  bool
	Length    bool
}

const (
	minUserLength = 6
	maxUserLength = 20
	minPswLength = 8
	maxPswLength = 20
)


// Checking for alphanumeric characters in username
func UsernameCorrect(u string) bool {
	for _, char := range u {
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			return false
		}
	}
	return true

}

// Checking for username len
func UsernameLen(u string) bool {
	if minUserLength <= len(u) && len(u) <= maxUserLength {
		return true
	}
	return false
}

// Password special conditions
func PswdConditions(p string) Password {
	var pw Password
	for _, char := range p {
		switch {
		case unicode.IsLower(char):
			pw.Lowercase = true
		case unicode.IsUpper(char):
			pw.Uppercase = true
		case unicode.IsNumber(char):
			pw.Number = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pw.Special = true
		case unicode.IsSpace(int32(char)):
			pw.NoSpaces = false
		}
	}

	if minPswLength <= len(p) && len(p) <= maxPswLength {
		pw.Length = true
	}
	return pw
}

// Checking if email follows email standards
func IsValidEmail(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(e)
}
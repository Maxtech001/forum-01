package registration

import "unicode"

type Password struct {
	Lowercase bool
	Uppercase bool
	Number    bool
	Special   bool
	NoSpaces  bool
	Length    bool
}


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

	if 5 <= len(u) && len(u) <= 50 {
		return true
	}
	return false

}

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

	if 11 < len(p) && len(p) < 60 {
		pw.Length = true
	}

	return pw
}
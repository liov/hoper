package prm

import "regexp"

const (
	Phone = iota
	Mail
	Unknown
)

func PhoneOrMail(input string) int {
	emailMatch, _ := regexp.MatchString(`^([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+.[a-zA-Z]{2,3}$`, input)
	if emailMatch {
		return Mail
	} else {
		phoneMatch, _ := regexp.MatchString(`^1[0-9]{10}$`, input)
		if phoneMatch {
			return Phone
		}
	}
	return Unknown
}

package common

import "regexp"

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
func IsValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^\+?[0-9]{7,15}$`)
	return re.MatchString(phone)
}
func IsValidOtpNumber(phone string) bool {
	re := regexp.MustCompile(`^\d{4}$`)
	return re.MatchString(phone)
}

package validator

func ValidateEmail(password string) bool {
	return len(password) > 3
}

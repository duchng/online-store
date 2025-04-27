package password

func Validate(password string) bool {
	if len(password) != 6 {
		return false
	}

	for _, c := range password {
		if c != rune(password[0]) {
			return true
		}
	}

	return false
}

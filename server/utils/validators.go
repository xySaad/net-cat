package utils

func ValidUsername(name string) bool {
	if len(name) < 3 || len(name) > 12 {
		return false
	}
	for _, char := range name {
		if (!(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'Z')) && !(char >= '0' && char <= '9') {
			return false
		}
	}
	return true
}

func ValidMsg(message []byte) bool {
	if len(message) == 0 {
		return false
	}
	for _, char := range string(message) {
		if char == 10 {
			continue
		}
		if char < 32 {
			return false
		}
	}
	return true
}

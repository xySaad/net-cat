package utils

func validUsername(name string) bool {
	if len(name) < 3 || len(name) > 12 {
		return false
	}
	for _, char := range name {
		if char == 27 {
			return false
		}
		if (!(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'z')) && !(char >= '0' && char <= '9') {
			return false
		}
	}
	return true
}

func validMsg(message []byte) bool {
	if len(message) == 0 {
		return false
	}
	for _, char := range string(message) {
		if (127 < char || char < 32) && char != 10 {
			return false
		}
	}
	return true
}

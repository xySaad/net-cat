package utils

func ValidName(name string) uint8 {
	if len(name) < 3 {
		return 1
	}
	if len(name) > 12 {
		return 2
	}
	for _, char := range name {
		if (!(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'Z')) && !(char >= '0' && char <= '9') && char != '_' {
			return 3
		}
	}
	return 0
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

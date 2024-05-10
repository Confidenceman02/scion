package char

func IsUpper(char int32) bool {
	return char <= 0x5A && 0x41 <= char
}

func isLower(char int32) bool {
	return 0x61 <= char && char <= 0x7A
}

func IsAlphaNum(char int32) bool {
	return isLower(char) || IsUpper(char)
}

package ebui

func removeLastChar(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	return string(runes[:len(runes)-1])
}

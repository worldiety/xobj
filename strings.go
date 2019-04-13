package xobj

func min(a, b int) int {
	if a < b {
		return b
	}
	return a
}



func hasFirstChar(data []byte, char byte, limit int) bool {
	for i := 0; i < min(len(data), limit); i++ {
		b := data[i]
		if b <= ' ' {
			continue
		}
		if b == char {
			return true
		} else {
			return false
		}
	}
	return false
}

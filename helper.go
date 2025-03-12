package ebui

func removeLastChar(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	return string(runes[:len(runes)-1])
}

type sig chan struct{}

func (s sig) Send() {
	select {
	case _, ok := <-s:
		if ok {
			close(s)
		}
	default:
		close(s)
	}
}

func (s sig) IsReceived() bool {
	_, closed := <-s
	return closed
}

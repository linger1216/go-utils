package pattern

type Stack []string

func (s Stack) pop() string {
	size := len(s)
	e := s[size-1]
	s = s[:size-1]
	return e
}

func (s Stack) push(e string) Stack {
	return append(s, e)
}

func (s Stack) peek() string {
	return s[len(s)-1]
}

func (s Stack) empty() bool {
	return len(s) == 0
}

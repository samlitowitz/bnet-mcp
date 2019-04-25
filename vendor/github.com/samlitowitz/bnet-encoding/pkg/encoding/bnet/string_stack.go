package bnet

type stringStack []string

func (s stringStack) Push(i string) stringStack {
	return append(s, i)
}

func (s stringStack) Pop() (stringStack, string) {
	l := len(s) - 1
	if l == -1 {
		panic("attempting to pop from an empty stack")
	}
	return s[:l], s[l]
}

func (s stringStack) Empty() bool {
	return len(s) < 1
}

func (s stringStack) Shift() (stringStack, string) {
	l := len(s) - 1
	if l == -1 {
		panic("attempting to pop from an empty stack")
	}
	return s[1:], s[0]
}

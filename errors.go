package bitbuf

type InsufficientError struct{}

func (e InsufficientError) Error() string {
	return "Insufficient"
}

type OverflowError struct{}

func (e OverflowError) Error() string {
	return "Overflow"
}

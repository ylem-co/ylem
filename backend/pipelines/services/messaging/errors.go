package messaging

type ErrorNonRepeatable struct {
	message string
}

func (e ErrorNonRepeatable) Error() string {
	return e.message
}

type ErrorRepeatable struct {
	message string
}

func (e ErrorRepeatable) Error() string {
	return e.message
}

func NewErrorNonRepeatable(message string) ErrorNonRepeatable {
	return ErrorNonRepeatable{
		message: message,
	}
}
func NewErrorRepeatable(message string) ErrorRepeatable {
	return ErrorRepeatable{
		message: message,
	}
}

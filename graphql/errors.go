package graphql

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrUserAlreadyExists = &Error{Code: "user_already_exists", Message: "User already exists"}
)
package utils

type Error struct {
	ErrCode int
	ErrMsg  string
}

func NewError(code int, msg string) *Error {
	return &Error{ErrCode: code, ErrMsg: msg}
}
func (err *Error) Error() string {
	return err.ErrMsg
}

func NewErrorDefault(msg string) *Error {
	return &Error{ErrCode: 500, ErrMsg: msg}
}

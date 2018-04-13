package cerror

import (
	"fmt"
)

type CodeError struct {
	Code     int
	InnerErr error
}

func NewCodeError(code int, err error) CodeError {
	return CodeError{code, err}
}

func (e CodeError) Error() string {
	return fmt.Sprintf("CodeError(%d):%s", e.Code, e.InnerErr.Error())
}

const (
	DB_ERROR = 121
)

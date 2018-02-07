package hdata

type Error struct {
	Err  error
	Code int
}

func (ar Error) Error() string {
	return ar.Err.Error()
}

func (ar Error) OK() bool {
	return ar.Err == nil
}

func NewError(code int, err error) Error {
	return Error{Err: err, Code: code}
}

var noError = NewError(0, nil)

func NoError() Error {
	return noError
}

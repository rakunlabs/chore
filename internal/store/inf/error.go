package inf

type ErrOperation struct {
	Err  error
	Code int
}

func (e ErrOperation) GetCode() int {
	return e.Code
}

func (e ErrOperation) Error() string {
	return e.Err.Error()
}

package app_errors

func NewErrNotFound(err error) ErrNotFound {
	return ErrNotFound{
		cause: err,
	}
}

type ErrNotFound struct {
	cause error
}

func (e ErrNotFound) Error() string {
	return "The requested resource could not be found"
}

func (e ErrNotFound) Is(target error) bool { return target.Error() == e.Error() }

func (e ErrNotFound) Unwrap() error {
	return e.cause
}
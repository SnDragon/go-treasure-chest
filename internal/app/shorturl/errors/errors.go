package errors

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Err  error
	Code int
}

func (s *StatusError) Error() string {
	return s.Err.Error()
}

func (s *StatusError) Status() int {
	return s.Code
}

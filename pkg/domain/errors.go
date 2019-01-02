package domain

// common errors used across packages

type notFoundErr struct {
	message string
}

func (nfe *notFoundErr) Error() string {
	return nfe.message
}

func NewNotFoundErr(missing string) *notFoundErr {
	return &notFoundErr{message: "did not find " + missing}
}

func IsNotFoundErr(err error) bool {
	_, ok := err.(*notFoundErr)
	return ok
}

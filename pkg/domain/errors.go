package domain

// common errors used across packages

type NotFoundErr struct {
	message string
}

type NotDirectMessageErr struct {
}

func (ndm *NotDirectMessageErr) Error() string {
	return "this command can only be done via a direct message"
}

func (nfe *NotFoundErr) Error() string {
	return nfe.message
}

func NewNotFoundErr(missing string) *NotFoundErr {
	return &NotFoundErr{message: "did not find " + missing}
}

func IsNotFoundErr(err error) bool {
	_, ok := err.(*NotFoundErr)
	return ok
}

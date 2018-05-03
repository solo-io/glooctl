package virtualservice

// a special kind of error returned by get default virtual service
type notExistsErr struct {
	message string
}

func (err notExistsErr) Error() string {
	return err.message
}

func NewNotExistsErr(msg string) notExistsErr {
	return notExistsErr{message: msg}
}

func IsNotExists(err error) bool {
	_, ok := err.(notExistsErr)
	return ok
}

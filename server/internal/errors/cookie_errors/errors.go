package cookie_errors

import "errors"

type ExpiredCookie struct {
	Err error
}

type ErrUnexistingCookie struct {
	Err error
}

func (e ExpiredCookie) Error() string {
	return e.Err.Error()
}

func (e ErrUnexistingCookie) Error() string {
	return e.Err.Error()
}

func NewExpiredCookie() ExpiredCookie {
	return ExpiredCookie{Err: errors.New("Cookie was deleted because of experation")}
}

func NewErrUnexistingCookie(name string) ErrUnexistingCookie {
	return ErrUnexistingCookie{Err: errors.New("Cookie with name " + name + " does not corresponds with any User")}
}

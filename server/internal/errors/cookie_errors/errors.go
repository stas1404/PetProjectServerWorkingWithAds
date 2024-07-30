package cookie_errors

import "errors"

type ExpiredCookie error

type ErrUnexistingCookie error

func NewExpiredCookie() ExpiredCookie {
	return errors.New("Cookie was deleted because of experation")
}

func NewErrUnexistingCookie(name string) ErrUnexistingCookie {
	return errors.New("Cookie with name " + name + " does not corresponds with any User")
}

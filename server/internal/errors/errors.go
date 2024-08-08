package errors

import (
	"errors"
	"strconv"
)

type PermissionDenied struct {
	Err error
}

type ErrUnexistingAd struct {
	Err error
}

type ErrUnexistingUser struct {
	Err error
}

type ErrBadUser struct {
	Err error
}

func (p PermissionDenied) Error() string {
	return p.Err.Error()
}

func (p ErrUnexistingAd) Error() string {
	return p.Err.Error()
}

func (p ErrUnexistingUser) Error() string {
	return p.Err.Error()
}

func (p ErrBadUser) Error() string {
	return p.Err.Error()
}

func NewPermissionDenied(UserID int64) PermissionDenied {
	return PermissionDenied{
		Err: errors.New("user with id " +
			strconv.FormatInt(UserID, 10) +
			" can not modify this ad as he is not its author"),
	}
}

func NewErrUnexistingAd(id int64) ErrUnexistingAd {
	return ErrUnexistingAd{Err: errors.New("Ad with id " + strconv.FormatInt(id, 10) + " does not exits")}
}

func NewErrUnexistingUser(id int64) ErrUnexistingUser {
	return ErrUnexistingUser{Err: errors.New("User with id " + strconv.FormatInt(id, 10) + " is not exists")}
}

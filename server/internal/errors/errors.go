package errors

import (
	"errors"
	"strconv"
)

type PermissionDenied struct {
	Err error
}

type ErrUnexistingAd error

type ErrUnexistingUser error

type ErrWrongUserID error

type ErrBadUser error

func (p PermissionDenied) Error() string {
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
	return errors.New("Ad with id " + strconv.FormatInt(id, 10) + " does not exits")
}

func NewErrUnexistingUser(id int64) ErrUnexistingUser {
	return errors.New("User with id " + strconv.FormatInt(id, 10) + " is not exists")
}

func NewErrWrongUserID(id int64) ErrWrongUserID {
	return errors.New("User with id " + strconv.FormatInt(id, 10) + " can not " +
		"modify this ad as he is not it's owner")
}

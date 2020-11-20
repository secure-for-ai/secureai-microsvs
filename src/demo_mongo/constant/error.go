package constant

import "errors"

/* Common Error */
var (
	ErrDatabase           = errors.New("ErrDatabase")
	ErrParamEmpty         = errors.New("ErrParamEmpty")
	ErrParamIDFormatWrong = errors.New("ErrParamIDFormatWrong")
)

/* User Info Error */
var (
	ErrAccountExist    = errors.New("ErrAccountExist")
	ErrAccountNotExist = errors.New("ErrAccountNotExist")
	ErrSession         = errors.New("ErrSession")
	ErrAccountNotLogin = errors.New("ErrAccountNotLogin")
)

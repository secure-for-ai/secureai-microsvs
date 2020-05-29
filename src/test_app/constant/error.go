package constant

import "errors"

/* Common Error */
var (
	ErrParamEmpty = errors.New("ErrParamEmpty")
)

/* User Info Error */
var (
	ErrAccountExist    = errors.New("ErrAccountExist")
	ErrAccountNotExist = errors.New("ErrAccountNotExist")
)

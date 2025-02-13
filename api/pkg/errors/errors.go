package haberrors

import "errors"

var (
	ErrBadRequest           = errors.New("bad request")
	ErrForbidden            = errors.New("forbidden")
	ErrRuntime              = errors.New("runtime error")
	ErrNotFound             = errors.New("not found")
	ErrValidation           = errors.New("validation error")
	ErrPreconditionFailed   = errors.New("precondition failed")
	ErrPreconditionRequired = errors.New("precondition required")
	ErrTimeout              = errors.New("timeout")
	ErrUnauthorized         = errors.New("unauthorized")
)

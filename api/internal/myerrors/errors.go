package myerrors

import "errors"

var (
	ErrBadRequest           = errors.New("alice: bad request")
	ErrForbidden            = errors.New("alice: forbidden")
	ErrRuntime              = errors.New("alice: runtime error")
	ErrNotFound             = errors.New("alice: not found")
	ErrValidation           = errors.New("alice: validation error")
	ErrPreconditionFailed   = errors.New("alice: precondition failed")
	ErrPreconditionRequired = errors.New("alice: precondition required")
	ErrTimeout              = errors.New("alice: timeout")
	ErrUnauthorized         = errors.New("alice: unauthorized")
)

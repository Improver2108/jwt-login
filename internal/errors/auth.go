package customerrors

import "errors"

var (
	ErrPasswordsMismatch    = errors.New("passwords dont match")
	ErrPasswordHash         = errors.New("cannot hash password")
	ErrIssuingTokens        = errors.New("errors issuing token")
	ErrTokenPersist         = errors.New("errors in persistin tokens")
	ErrWrongPassword        = errors.New("wrong password error")
	ErrCredentialValidation = errors.New("error validating credentials")
)

package constant

const (
	// general
	DataAlreadyExist   = "data already exists"
	DataNotFound       = "data not found"
	RecordNotFound     = "record not found"
	FileSizeIsTooLarge = "file size should not exceed 200 KB"
	InvalidDateFormat  = "invalid date format"

	// auth
	AlreadyVerified            = "the account has already verified"
	InvalidEmailOrPassword     = "email or password is incorrect"
	InvalidPassword            = "password must contain a combination of uppercase, lowercase, number, and symbol"
	InvalidPasswordResetLink   = "invalid password reset link"
	InvalidActivationLink      = "invalid activation link"
	IncorrectPassword          = "password is incorrect"
	ConfirmNewPasswordMismatch = "please ensure that the new password and confirm password fields match exactly"
	ConfirmPasswordMismatch    = "please ensure that password and confirm password fields match exactly"
	RequestProhibited          = "request is prohibited"

	InvalidActiveValue = "invalid value for 'active', it must be a boolean"
	EmailAlreadyExists = "email already exists"
)

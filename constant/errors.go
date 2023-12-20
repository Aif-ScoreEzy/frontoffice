package constant

const (
	// general
	DataAlreadyExist   = "data already exists"
	DataNotFound       = "data not found"
	FileSizeIsTooLarge = "file size should not exceed 200 KB"
	InvalidDateFormat  = "invalid date format"
	RecordNotFound     = "record not found"

	// auth
	AlreadyVerified            = "the account has already verified"
	ActivationTokenExpired     = "user activation token has expired"
	ConfirmNewPasswordMismatch = "please ensure that the new password and confirm password fields match exactly"
	ConfirmPasswordMismatch    = "please ensure that password and confirm password fields match exactly"
	InvalidEmailOrPassword     = "email or password is incorrect"
	InvalidPassword            = "password must contain a combination of uppercase, lowercase, number, and symbol"
	InvalidPasswordResetLink   = "invalid password reset link"
	InvalidActivationLink      = "invalid activation link"
	IncorrectPassword          = "password is incorrect"
	RequestProhibited          = "request is prohibited"
	TokenExpired               = "Token is expired"

	DuplicateGrading = "duplicate grading"

	EmailAlreadyExists = "email already exists"
	InvalidImageFile   = "invalid image file"
	InvalidStatusValue = "invalid value for 'status'"
	SendEmailFailed    = "send email failed"
)

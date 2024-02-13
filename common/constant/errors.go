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
	UnverifiedUser             = "please check your email, you need to verify your email address before you can reset your password"
	TokenExpired               = "Token is expired"

	//grading
	DuplicateGrading       = "duplicate grading"
	FieldGradingLabelEmpty = "field grading label is required"
	FieldMinGradeEmpty     = "field min grade is required"
	FieldMaxGradeEmpty     = "field max grade is required"
	FieldGradingValueEmpty = "field grading value is required"

	// gen-retail
	InvalidDocumentFile    = "invalid document file"
	ErrorGettingFile       = "error getting file"
	ErrorOpeningFile       = "error opening file"
	ErrorReadingCSV        = "error reading CSV file"
	HeaderTemplateNotValid = "header template is not valid"
	ErrorReadingCSVRecords = "error reading CSV records"
	ErrorUploadDataCSV     = "error upload data CSV file"

	//parameter settings
	ParamSettingIsNotSet = "parameter settings is not set"

	EmailAlreadyExists = "email already exists"
	InvalidImageFile   = "invalid image file"
	InvalidStatusValue = "invalid value for 'status'"
	SendEmailFailed    = "send email failed"
)

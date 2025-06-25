package loanrecordchecker

type loanRecordCheckerRequest struct {
	Name  string `json:"name" validate:"required~Name cannot be empty"`
	Nik   string `json:"nik" validate:"required~NIK cannot be empty., numeric~ID Card No is only number, length(16)~ID Card No must be 16 digit number."`
	Phone string `json:"phone_number" validate:"required~Phone Number cannot be empty, indophone, min(9)"`
}

type dataLoanRecord struct {
	Remarks string `json:"remarks"`
	Status  string `json:"status"`
}

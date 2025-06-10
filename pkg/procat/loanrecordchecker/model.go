package loanrecordchecker

type LoanRecordCheckerRequest struct {
	Name  string `json:"name" validate:"required~Name cannot be empty"`
	Nik   string `json:"nik" validate:"required~NIK cannot be empty., numeric~ID Card No is only number, length(16)~ID Card No must be 16 digit number."`
	Phone string `json:"phone_number" validate:"required~Phone Number cannot be empty, indophone, min(9)"`
}

type LoanRecordCheckerRawResponse struct {
	Success         bool        `json:"success"`
	Data            interface{} `json:"data"`
	PricingStrategy interface{} `json:"pricing_strategy"`
	TransactionId   interface{} `json:"transaction_id"`
	DateTime        interface{} `json:"datetime"`
	Message         string      `json:"message"`
	StatusCode      int         `json:"status_code"`
}

type LoanRecordCheckerResponse struct {
	Data            interface{} `json:"data"`
	PricingStrategy interface{} `json:"pricing_strategy"`
	TransactionID   interface{} `json:"transaction_id"`
	Datetime        interface{} `json:"datetime"`
}

type loanRecordCheckerFilter struct {
	Page        string
	Size        string
	Offset      string
	StartDate   string
	EndDate     string
	ProductSlug string
	MemberId    string
	CompanyId   string
	TierLevel   string
	Keyword     string
}

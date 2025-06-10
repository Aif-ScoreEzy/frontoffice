package multipleloan

type MultipleLoanRequest struct {
	Nik   string `json:"nik" validate:"required~NIK cannot be empty., numeric~ID Card No is only number, length(16)~ID Card No must be 16 digit number."`
	Phone string `json:"phone_number" validate:"required~Phone Number cannot be empty, indophone, min(9)"`
}

type dataMultipleLoanResponse struct {
	QueryCount uint `json:"query_count"`
}

type multipleLoanFilter struct {
	Page        string
	Size        string
	Offset      string
	StartDate   string
	EndDate     string
	JobId       string
	ProductSlug string
	MemberId    string
	CompanyId   string
	TierLevel   string
}

type MultipleLoanRawResponse struct {
	Success         bool        `json:"success"`
	Data            interface{} `json:"data"`
	PricingStrategy interface{} `json:"pricing_strategy"`
	TransactionId   interface{} `json:"transaction_id"`
	DateTime        interface{} `json:"datetime"`
	Message         string      `json:"message"`
	StatusCode      int         `json:"status_code"`
}

type MultipleLoanResponse struct {
	Data            interface{} `json:"data"`
	PricingStrategy interface{} `json:"pricing_strategy"`
	TransactionID   interface{} `json:"transaction_id"`
	Datetime        interface{} `json:"datetime"`
}

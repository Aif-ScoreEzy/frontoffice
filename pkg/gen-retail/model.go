package genretail

type GenRetailRequest struct {
	LoanNo   string `json:"loan_no"`
	Name     string `json:"name"`
	IDCardNo string `json:"id_card_no"`
	PhoneNo  string `json:"phone_no"`
}

type GenRetailV3ModelResponse struct {
	Message      string                 `json:"message"`
	ErrorMessage string                 `json:"error_message"`
	Success      bool                   `json:"success"`
	Data         *GenRetailV3DataClient `json:"data"`
	StatusCode   int                    `json:"status_code"`
}

type GenRetailV3DataClient struct {
	TransactionID        string  `json:"transaction_id"`
	Name                 string  `json:"name"`
	IDCardNo             string  `json:"id_card_no"`
	PhoneNo              string  `json:"phone_no"`
	LoanNo               string  `json:"loan_no"`
	ProbabilityToDefault float64 `json:"probability_to_default"`
	Grade                string  `json:"grade"`
	Date                 string  `json:"date"`
}

type GenRetailV3ClientReturnSuccess struct {
	Message string                `json:"message"`
	Success bool                  `json:"success"`
	Data    GenRetailV3DataClient `json:"data"`
}

type GenRetailV3ClientReturnError struct {
	Message      string                 `json:"message"`
	ErrorMessage string                 `json:"error_message"`
	Data         *GenRetailV3DataClient `json:"data"`
}

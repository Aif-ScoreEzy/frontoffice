package multipleloan

type multipleLoanRequest struct {
	Nik   string `json:"nik" validate:"required~NIK cannot be empty., numeric~ID Card No is only number, length(16)~ID Card No must be 16 digit number."`
	Phone string `json:"phone_number" validate:"required~Phone Number cannot be empty, indophone, min(9)"`
}

type dataMultipleLoanResponse struct {
	QueryCount uint `json:"query_count"`
}

package job

import (
	"time"
)

type logTransProductCatalog struct {
	MemberID        uint           `json:"member_id"`
	CompanyID       uint           `json:"company_id"`
	JobID           uint           `json:"job_id"`
	ProductID       uint           `json:"product_id"`
	Status          string         `json:"status"`
	Message         *string        `json:"message"`
	Input           *logTransInput `json:"input"`
	Data            *logTransData  `json:"data"`
	PricingStrategy string         `json:"pricing_strategy"`
	TransactionId   string         `json:"transaction_id"`
	DateTime        string         `json:"datetime"`
}

type logTransData struct {
	Remarks string `json:"remarks"`
	Status  string `json:"status"`
}

type logTransInput struct {
	Name        string `json:"name"`
	NIK         string `json:"nik"`
	PhoneNumber string `json:"phone_number"`
}

type jobDetailResponse struct {
	TotalData                  int64                     `json:"total_data"`
	TotalDataPercentageSuccess int64                     `json:"total_data_percentage_success"`
	TotalDataPercentageFail    int64                     `json:"total_data_percentage_fail"`
	TotalDataPercentageError   int64                     `json:"total_data_percentage_error"`
	JobDetails                 []*logTransProductCatalog `json:"job_details"`
}

type CreateJobRequest struct {
	ProductId uint   `json:"product_id" validate:"required~Field product id is required"`
	MemberId  string `json:"member_id" validate:"required~Field member id is required"`
	CompanyId string `json:"company_id" validate:"required~Field company id is required"`
	Total     int    `json:"total" validate:"required~Field total is required"`
}

type UpdateJobRequest struct {
	SuccessCount *uint      `json:"success_count"`
	Status       *string    `json:"status"`
	EndAt        *time.Time `json:"end_at"`
}

type createJobRespData struct {
	JobId uint `json:"id"`
}

type logFilter struct {
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

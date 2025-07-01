package job

import (
	"time"
)

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

type createJobDataResponse struct {
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

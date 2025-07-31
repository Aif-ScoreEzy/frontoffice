package oldphonelivestatus

import (
	"front-office/pkg/core/company"
	"front-office/pkg/core/member"
	"time"
)

type mstPhoneLiveStatusJob struct {
	Id           uint               `json:"id"`
	Total        int                `json:"total"`
	SuccessCount int                `json:"success_count"`
	Status       string             `json:"status"`
	MemberId     uint               `json:"member_id"`
	Member       member.MstMember   `json:"-"`
	CompanyId    uint               `json:"company_id"`
	Company      company.MstCompany `json:"-"`
	CreatedAt    time.Time          `json:"start_time"`
	EndAt        *time.Time         `json:"end_time"`
}

type mstPhoneLiveStatusJobDetail struct {
	Id               uint      `json:"id"`
	MemberId         uint      `json:"member_id"`
	CompanyId        uint      `json:"company_id"`
	JobId            uint      `json:"job_id"`
	PhoneNumber      string    `json:"phone_number" validate:"required~phone number is required, min(10)~phone number must be at least 10 characters, indophone~invalid number"`
	InProgress       bool      `json:"in_progess"`
	Sequence         int       `json:"sequence"`
	Status           string    `json:"status"`
	Message          *string   `json:"message"`
	SubscriberStatus string    `json:"subscriber_status"`
	DeviceStatus     string    `json:"device_status"`
	PhoneType        string    `json:"phone_type"`
	Operator         string    `json:"operator"`
	PricingStrategy  string    `json:"pricing_strategy"`
	TransactionId    string    `json:"transaction_id"`
	CreatedAt        time.Time `json:"created_at"`
}

type phoneLiveStatusRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required~phone number is required, min(10)~phone number must be at least 10 characters, indophone~invalid number"`
	TrxId       string `json:"trx_id"`
}

type phoneLiveStatusRespData struct {
	LiveStatus string      `json:"live_status"`
	PhoneType  string      `json:"phone_type"`
	Operator   string      `json:"operator"`
	Errors     []errorData `json:"errors"`
}

type errorData struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type phoneLiveStatusFilter struct {
	Page      string
	Size      string
	Offset    string
	StartDate string
	EndDate   string
	JobId     string
	MemberId  string
	CompanyId string
	TierLevel string
	Keyword   string
}

type jobListRespData struct {
	Jobs      []*mstPhoneLiveStatusJob `json:"jobs"`
	TotalData int                      `json:"total_data"`
}

type jobDetailRespData struct {
	TotalData                  int64                          `json:"total_data"`
	TotalDataPercentageSuccess int64                          `json:"total_data_percentage_success"`
	TotalDataPercentageFail    int64                          `json:"total_data_percentage_fail"`
	TotalDataPercentageError   int64                          `json:"total_data_percentage_error"`
	SubsActive                 int64                          `json:"subs_active"`
	DevReachable               int64                          `json:"dev_reachable"`
	DevUnreachable             int64                          `json:"dev_unreachable"`
	DevUnavailable             int64                          `json:"dev_unavailable"`
	JobDetails                 []*mstPhoneLiveStatusJobDetail `json:"job_details"`
}

type jobsSummaryRespData struct {
	TotalData        int64 `json:"total_data"`
	TotalDataSuccess int64 `json:"total_data_percentage_success"`
	TotalDataFail    int64 `json:"total_data_percentage_fail"`
	TotalDataError   int64 `json:"total_data_percentage_error"`
	SubscriberActive int64 `json:"subs_active"`
	DeviceReachable  int64 `json:"dev_reachable"`
	Mobile           int64 `json:"mobile"`
	FixedLine        int64 `json:"fixed_line"`
}

type createJobRequest struct {
	MemberId                string                    `json:"member_id"`
	CompanyId               string                    `json:"company_id"`
	PhoneLiveStatusRequests []*phoneLiveStatusRequest `json:"requests"`
}

type createJobRespData struct {
	JobId uint `json:"job_id"`
}

type getSuccessCountRespData struct {
	SuccessCount uint `json:"success_count"`
}

type updateJobRequest struct {
	SuccessCount *int       `json:"success_count"`
	Status       *string    `json:"status"`
	EndAt        *time.Time `json:"end_at"`
}

type updateJobDetailRequest struct {
	InProgress       *bool      `json:"in_progress"`
	Sequence         *int       `json:"sequence"`
	Status           *string    `json:"status"`
	Message          *string    `json:"message"`
	SubscriberStatus *string    `json:"subscriber_status"`
	DeviceStatus     *string    `json:"device_status"`
	PhoneType        *string    `json:"phone_type"`
	Operator         *string    `json:"operator"`
	PricingStrategy  *string    `json:"pricing_strategy"`
	TransactionId    *string    `json:"transaction_id"`
	EndAt            *time.Time `json:"end_time"`
}

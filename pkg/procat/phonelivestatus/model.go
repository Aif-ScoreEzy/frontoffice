package phonelivestatus

import (
	"front-office/pkg/core/company"
	"front-office/pkg/core/member"
	"time"
)

type MstPhoneLiveStatusJob struct {
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

type MstPhoneLiveStatusJobDetail struct {
	Id               uint                  `json:"id"`
	MemberId         uint                  `json:"member_id"`
	Member           member.MstMember      `json:"-"`
	CompanyId        uint                  `json:"company_id"`
	Company          company.MstCompany    `json:"-"`
	JobId            uint                  `json:"job_id"`
	Job              MstPhoneLiveStatusJob `json:"-"`
	PhoneNumber      string                `json:"phone_number"`
	SubscriberStatus string                `json:"subscriber_status"`
	DeviceStatus     string                `json:"device_status"`
	InProgress       bool                  `json:"in_progess"`
	Sequence         int                   `json:"sequence"`
	Status           string                `json:"status"`
	Message          *string               `json:"message"`
	CreatedAt        time.Time             `json:"-"`
}

type PhoneLiveStatusRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type PhoneLiveStatusFilter struct {
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

type APIResponse[T any] struct {
	Success    bool   `json:"success"`
	Data       T      `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

type JobListResponse struct {
	Jobs      []*MstPhoneLiveStatusJob `json:"jobs"`
	TotalData int                      `json:"total_data"`
}

type JobDetailsResponse struct {
	TotalData                  int64                          `json:"total_data"`
	TotalDataPercentageSuccess int64                          `json:"total_data_percentage_success"`
	TotalDataPercentageFail    int64                          `json:"total_data_percentage_fail"`
	TotalDataPercentageError   int64                          `json:"total_data_percentage_error"`
	SubsActive                 int64                          `json:"subs_active"`
	DevReachable               int64                          `json:"dev_reachable"`
	DevUnreachable             int64                          `json:"dev_unreachable"`
	DevUnavailable             int64                          `json:"dev_unavailable"`
	JobDetails                 []*MstPhoneLiveStatusJobDetail `json:"job_details"`
}

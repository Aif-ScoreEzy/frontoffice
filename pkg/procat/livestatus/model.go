package livestatus

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Job struct {
	Id        uint       `json:"id"`
	Total     int        `json:"total"`
	Success   int        `json:"success"`
	Status    string     `json:"status"`
	UserId    string     `json:"user_id"`
	CompanyId string     `json:"company_id"`
	CreatedAt time.Time  `gorm:"not null;default:current_timestamp" json:"start_time"`
	EndAt     *time.Time `json:"end_time"`
}

type JobDetail struct {
	Id               uint      `json:"id"`
	UserId           string    `json:"user_id"`
	CompanyId        string    `json:"company_id"`
	JobId            uint      `json:"job_id"`
	PhoneNumber      string    `json:"phone_number" validate:"required~phone number is required, min(10)~phone number must be at least 10 characters, indophone~invalid number"`
	SubscriberStatus string    `json:"subscriber_status"`
	DeviceStatus     string    `json:"device_status"`
	OnProcess        bool      `gorm:"not null" json:"on_process"`
	Sequence         int       `json:"sequence"`
	Status           string    `json:"status"`
	Message          *string   `json:"message"`
	Data             *JSONB    `gorm:"type:jsonb" json:"data"`
	CreatedAt        time.Time `gorm:"not null;default:current_timestamp" json:"-"`
}

type UpdateJobRequest struct {
	Total  *int       `json:"total"`
	Status *string    `json:"status"`
	EndAt  *time.Time `json:"end_at"`
}

type UpdateJobDetailRequest struct {
	Id               uint      `json:"id"`
	JobId            uint      `json:"job_id"`
	PhoneNumber      string    `json:"phone_number"`
	SubscriberStatus string    `json:"subscriber_status"`
	DeviceStatus     string    `json:"device_status"`
	OnProcess        bool      `json:"on_process"`
	Sequence         int       `json:"sequence"`
	Status           string    `json:"status"`
	Data             *JSONB    `json:"data"`
	CreatedAt        time.Time `json:"-"`
}

type LiveStatusRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required~phone number is required, min(10)~phone number must be at least 10 characters, indophone~invalid number"`
	TrxId       string `json:"trx_id"`
}

type LiveStatusResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
}

type JobSummaryResponse struct {
	TotalData        int64 `json:"total_data"`
	TotalDataSuccess int64 `json:"total_data_percentage_success"`
	TotalDataFail    int64 `json:"total_data_percentage_fail"`
	TotalDataError   int64 `json:"total_data_percentage_error"`
	SubscriberActive int64 `json:"subs_active"`
	DeviceReachable  int64 `json:"dev_reachable"`
	Mobile           int64 `json:"mobile"`
	FixedLine        int64 `json:"fixed_line"`
}

type GetJobsResponse struct {
	TotalData int64  `json:"total_data"`
	Jobs      []*Job `json:"jobs"`
}

type JobDetailResponse struct {
	TotalData         int64                   `json:"total_data"`
	TotalDataSuccess  int64                   `json:"total_data_percentage_success"`
	TotalDataFail     int64                   `json:"total_data_percentage_fail"`
	TotalDataError    int64                   `json:"total_data_percentage_error"`
	SubscriberActive  int64                   `json:"subs_active"`
	DeviceReachable   int64                   `json:"dev_reachable"`
	DeviceUnreachable int64                   `json:"dev_unreachable"`
	DeviceUnavailable int64                   `json:"dev_unavailable"`
	JobDetails        []*JobDetailQueryResult `json:"job_details"`
}

type JobDetailQueryResult struct {
	Id               uint   `json:"id"`
	JobId            uint   `json:"job_id"`
	PhoneNumber      string `json:"phone_number"`
	SubscriberStatus string `json:"subscriber_status"`
	DeviceStatus     string `json:"device_status"`
	Status           string `json:"status"`
	Operator         string `json:"operator"`
	PhoneType        string `json:"phone_type"`
	Message          string `json:"message"`
}

type ResponseSuccess struct {
	Success   int `json:"success"`
	TotalData int `json:"total_data"`
}

type JSONB map[string]interface{}

func (jsonField JSONB) Value() (driver.Value, error) {
	return json.Marshal(jsonField)
}

func (jsonField *JSONB) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(data, &jsonField)
}

func (j *JSONB) Encoded(value []byte) error {
	ab := json.Unmarshal(value, j)

	return ab
}

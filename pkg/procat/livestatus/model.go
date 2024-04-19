package livestatus

import "time"

type Job struct {
	ID        uint      `json:"id"`
	Total     int       `json:"total"`
	Success   int       `json:"success"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp" json:"-"`
}

type JobDetail struct {
	ID               uint      `json:"id"`
	JobID            uint      `json:"job_id"`
	PhoneNumber      string    `json:"phone_number"`
	SubscriberStatus string    `json:"subscriber_status"`
	DeviceStatus     string    `json:"device_status"`
	OnProcess        bool      `gorm:"not null;default:true" json:"on_process"`
	Sequence         int       `json:"sequence"`
	CreatedAt        time.Time `gorm:"not null;default:current_timestamp" json:"-"`
}

type UpdateJobDetailRequest struct {
	ID               uint      `json:"id"`
	JobID            uint      `json:"job_id"`
	PhoneNumber      string    `json:"phone_number"`
	SubscriberStatus string    `json:"subscriber_status"`
	DeviceStatus     string    `json:"device_status"`
	OnProcess        bool      `json:"on_process"`
	Sequence         int       `json:"sequence"`
	CreatedAt        time.Time `json:"-"`
}

type LiveStatusRequest struct {
	PhoneNumber string `json:"phone_number"`
	TrxID       string `json:"trx_id"`
}

type LiveStatusResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
}

type ResponseSuccess struct {
	Success   int `json:"success"`
	TotalData int `json:"total_data"`
}

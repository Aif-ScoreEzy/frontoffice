package livestatus

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Job struct {
	ID        uint      `json:"id"`
	Total     int       `json:"total"`
	Success   int       `json:"success"`
	Status    string    `json:"status"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp" json:"-"`
}

type JobDetail struct {
	ID               uint      `json:"id"`
	JobID            uint      `json:"job_id"`
	PhoneNumber      string    `json:"phone_number"`
	SubscriberStatus string    `json:"subscriber_status"`
	DeviceStatus     string    `json:"device_status"`
	OnProcess        bool      `gorm:"not null" json:"on_process"`
	Sequence         int       `json:"sequence"`
	Status           string    `json:"status"`
	Data             *JSONB    `gorm:"type:jsonb" json:"data"`
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
	Status           string    `json:"status"`
	Data             *JSONB    `json:"data"`
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

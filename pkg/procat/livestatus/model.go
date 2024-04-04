package livestatus

import "time"

type Job struct {
	ID        uint      `json:"id"`
	Total     int       `json:"total"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp" json:"-"`
}

type JobDetail struct {
	ID          uint      `json:"id"`
	JobID       uint      `json:"job_id"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `gorm:"not null;default:current_timestamp" json:"-"`
}

type FIFRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type FIFRequests struct {
	PhoneNumbers []FIFRequest `json:"phone_numbers"`
}

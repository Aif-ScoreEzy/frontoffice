package log

import (
	"time"

	"gorm.io/datatypes"
)

type LogTransScoreezy struct {
	LogTrxID             uint           `json:"log_trx_id" gorm:"primaryKey;autoIncrement"`
	TrxID                string         `json:"trx_id"`
	MemberID             uint           `json:"user_id"`
	CompanyID            uint           `json:"company_id"`
	IpClient             string         `json:"ip_client"`
	ProductID            uint           `json:"product_id"`
	Status               string         `json:"status"`  // Free or Pay
	Success              bool           `json:"success"` // true or false
	Message              string         `json:"message"`
	ProbabilityToDefault string         `json:"probability_to_default"`
	Grade                string         `json:"grade"`
	LoanNo               string         `json:"loan_no"`
	Data                 datatypes.JSON `json:"data" swaggertype:"object"`
	Duration             time.Duration  `json:"duration" format:"duration" example:"2h30m"`
	CreatedAt            time.Time      `json:"created_at" format:"date-time"`
}

type FetchLogTransResponse struct {
	Message    string        `json:"message"`
	Success    bool          `json:"success"`
	Data       *DataLogTrans `json:"data"`
	StatusCode int           `json:"-"`
}

type DataLogTrans struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Grade     string    `json:"grade"`
	CreatedAt time.Time `json:"created_at"`
}

type AifResponse struct {
	Success bool               `json:"success"`
	Data    []LogTransScoreezy `json:"data"`
	Message string             `json:"message"`
	Meta    any                `json:"meta,omitempty"`
	Status  bool               `json:"status,omitempty"`
}

type Meta struct {
	Message   string `json:"message"`
	Total     any    `json:"total,omitempty"`
	Page      any    `json:"page,omitempty"`
	TotalPage any    `json:"total_page,omitempty"`
	Visible   any    `json:"visible,omitempty"`
	StartData any    `json:"start_data,omitempty"`
	EndData   any    `json:"end_data,omitempty"`
	Size      any    `json:"size,omitempty"`
}

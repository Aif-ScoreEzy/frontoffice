package transaction

import (
	"front-office/pkg/core/member"
	"time"

	"gorm.io/datatypes"
)

type LogTransScoreezy struct {
	LogTrxId             uint             `json:"log_trx_id" gorm:"primaryKey;autoIncrement"`
	TrxId                string           `json:"trx_id"`
	MemberId             uint             `json:"user_id"`
	Member               member.MstMember `json:"member"`
	CompanyId            uint             `json:"company_id"`
	IpClient             string           `json:"ip_client"`
	ProductId            uint             `json:"product_id"`
	Status               string           `json:"status"`  // Free or Pay
	Success              bool             `json:"success"` // true or false
	Message              string           `json:"message"`
	ProbabilityToDefault string           `json:"probability_to_default"`
	Grade                string           `json:"grade"`
	LoanNo               string           `json:"loan_no"`
	Data                 datatypes.JSON   `json:"data" swaggertype:"object"`
	Duration             time.Duration    `json:"duration" format:"duration" example:"2h30m"`
	CreatedAt            time.Time        `json:"created_at" format:"date-time"`
}

type scoreezyLogResponse struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Grade     string    `json:"grade"`
	CreatedAt time.Time `json:"created_at"`
}

type getSuccessCountDataResponse struct {
	SuccessCount uint `json:"success_count"`
}

type UpdateTransRequest struct {
	Success *bool `json:"success"`
}

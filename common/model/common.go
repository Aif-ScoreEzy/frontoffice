package model

type AifResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
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

type ProCatAPIResponse[T any] struct {
	Success         bool        `json:"success"`
	Data            T           `json:"data"`
	Input           interface{} `json:"input"`
	Message         string      `json:"message"`
	StatusCode      int         `json:"-"`
	PricingStrategy interface{} `json:"pricing_strategy"`
	TransactionId   interface{} `json:"transaction_id"`
	Date            interface{} `json:"datetime"`
}

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

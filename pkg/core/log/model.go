package log

type AifResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message"`
	Meta    any    `json:"meta,omitempty"`
	Status  bool   `json:"status,omitempty"`
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

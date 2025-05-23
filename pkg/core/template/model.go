package template

type DownloadRequest struct {
	Category string `query:"category"` // Enum validation
	Filename string `query:"filename"`
}

type TemplateInfo struct {
	Category    string
	Files       []string
	Description string
}

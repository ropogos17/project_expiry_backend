package models

// Pagination struct for API responses
type Pagination struct {
	Page     int   `json:"page" form:"page" binding:"min=1"`
	PageSize int   `json:"page_size" form:"page_size" binding:"min=1,max=100"`
	Total    int64 `json:"total"`
	Pages    int   `json:"pages"`
}

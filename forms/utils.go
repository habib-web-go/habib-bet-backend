package forms

type PaginationMetadata struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	PageCount int `json:"page_count"`
}

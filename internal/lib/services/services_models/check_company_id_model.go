package services_models

type CompanyInfo struct {
	CompanyId   int64  `json:"company_id" validate:"required"`
	CompanyName string `json:"company_name" validate:"required"`
}

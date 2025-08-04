package services_models

type CreateBookingTypeDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	CompanyId   int64  `json:"company_id" validate:"required"`
	CompanyName string `json:"company_name" validate:"required"`
}

package model

type UpdateProduct struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Category    string  `json:"category" binding:"required"`
}

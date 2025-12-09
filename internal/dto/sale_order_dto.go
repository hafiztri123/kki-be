package dto

import "github.com/google/uuid"

type CreateSaleOrderRequest struct {
	CustomerName string  `json:"customer_name"`
	TotalAmount  float64 `json:"total_amount"`
	Status       string  `json:"status"`
}

type UpdateSaleOrderRequest struct {
	CustomerName string  `json:"customer_name"`
	TotalAmount  float64 `json:"total_amount"`
	Status       string  `json:"status"`
}

type SaleOrderResponse struct {
	ID           uuid.UUID `json:"id"`
	OrderNumber  string    `json:"order_number"`
	CustomerName string    `json:"customer_name"`
	TotalAmount  float64   `json:"total_amount"`
	Status       string    `json:"status"`
	CreatedBy    uuid.UUID `json:"created_by"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type PaginationRequest struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type PaginatedResponse struct {
	Data       any   `json:"data"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
}

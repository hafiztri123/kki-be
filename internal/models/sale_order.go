package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SaleOrder struct {
	ID          uuid.UUID    `json:"id"`
	OrderNumber string       `json:"order_number"`
	CustomerName string      `json:"customer_name"`
	TotalAmount float64      `json:"total_amount"`
	Status      string       `json:"status"`
	CreatedBy   uuid.UUID    `json:"created_by"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at,omitempty"`
}

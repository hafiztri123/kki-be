package repository

import (
	"context"
	"database/sql"
	"errors"

	apperror "github.com/hafiztri123/kki-be/internal/app_error"
	"github.com/hafiztri123/kki-be/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleOrderRepository struct {
	db *pgxpool.Pool
}

func NewSaleOrderRepository(db *pgxpool.Pool) *SaleOrderRepository {
	return &SaleOrderRepository{
		db: db,
	}
}

func (r *SaleOrderRepository) InsertSaleOrder(ctx context.Context, saleOrder *models.SaleOrder) error {
	query := `INSERT INTO sale_orders (id, order_number, customer_name, total_amount, status, created_by, created_at, updated_at, deleted_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(ctx, query,
		saleOrder.ID,
		saleOrder.OrderNumber,
		saleOrder.CustomerName,
		saleOrder.TotalAmount,
		saleOrder.Status,
		saleOrder.CreatedBy,
		saleOrder.CreatedAt,
		saleOrder.UpdatedAt,
		saleOrder.DeletedAt,
	)

	return err
}

func (r *SaleOrderRepository) GetSaleOrders(ctx context.Context, limit, offset int) ([]models.SaleOrder, int64, error) {
	var totalCount int64
	countQuery := `SELECT COUNT(*) FROM sale_orders WHERE deleted_at IS NULL`
	err := r.db.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, order_number, customer_name, total_amount, status, created_by, created_at, updated_at, deleted_at
			  FROM sale_orders
			  WHERE deleted_at IS NULL
			  ORDER BY created_at DESC
			  LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var saleOrders []models.SaleOrder
	for rows.Next() {
		var so models.SaleOrder
		err := rows.Scan(
			&so.ID,
			&so.OrderNumber,
			&so.CustomerName,
			&so.TotalAmount,
			&so.Status,
			&so.CreatedBy,
			&so.CreatedAt,
			&so.UpdatedAt,
			&so.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		saleOrders = append(saleOrders, so)
	}

	return saleOrders, totalCount, nil
}

func (r *SaleOrderRepository) GetSaleOrderByID(ctx context.Context, id uuid.UUID) (*models.SaleOrder, error) {
	query := `SELECT id, order_number, customer_name, total_amount, status, created_by, created_at, updated_at, deleted_at
			  FROM sale_orders
			  WHERE id = $1 AND deleted_at IS NULL`

	var so models.SaleOrder
	err := r.db.QueryRow(ctx, query, id).Scan(
		&so.ID,
		&so.OrderNumber,
		&so.CustomerName,
		&so.TotalAmount,
		&so.Status,
		&so.CreatedBy,
		&so.CreatedAt,
		&so.UpdatedAt,
		&so.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	return &so, nil
}

func (r *SaleOrderRepository) UpdateSaleOrder(ctx context.Context, saleOrder *models.SaleOrder) error {
	query := `UPDATE sale_orders
			  SET customer_name = $1, total_amount = $2, status = $3, updated_at = $4
			  WHERE id = $5 AND deleted_at IS NULL`

	result, err := r.db.Exec(ctx, query,
		saleOrder.CustomerName,
		saleOrder.TotalAmount,
		saleOrder.Status,
		saleOrder.UpdatedAt,
		saleOrder.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

func (r *SaleOrderRepository) DeleteSaleOrder(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE sale_orders
			  SET deleted_at = NOW()
			  WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

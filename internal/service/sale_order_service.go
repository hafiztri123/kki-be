package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hafiztri123/kki-be/internal/dto"
	"github.com/hafiztri123/kki-be/internal/models"
	"github.com/hafiztri123/kki-be/internal/repository"
)

type SaleOrderService struct {
	saleOrderRepo *repository.SaleOrderRepository
}

func NewSaleOrderService(saleOrderRepo *repository.SaleOrderRepository) *SaleOrderService {
	return &SaleOrderService{
		saleOrderRepo: saleOrderRepo,
	}
}

func (s *SaleOrderService) CreateSaleOrder(ctx context.Context, req *dto.CreateSaleOrderRequest, createdBy uuid.UUID) error {
	orderNumber := fmt.Sprintf("SO-%d", time.Now().Unix())

	saleOrder := &models.SaleOrder{
		ID:           uuid.New(),
		OrderNumber:  orderNumber,
		CustomerName: req.CustomerName,
		TotalAmount:  req.TotalAmount,
		Status:       req.Status,
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    sql.NullTime{},
	}

	return s.saleOrderRepo.InsertSaleOrder(ctx, saleOrder)
}

func (s *SaleOrderService) GetSaleOrders(ctx context.Context, limit, offset int) ([]dto.SaleOrderResponse, int64, error) {
	saleOrders, totalCount, err := s.saleOrderRepo.GetSaleOrders(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.SaleOrderResponse, 0, len(saleOrders))
	for _, so := range saleOrders {
		responses = append(responses, dto.SaleOrderResponse{
			ID:           so.ID,
			OrderNumber:  so.OrderNumber,
			CustomerName: so.CustomerName,
			TotalAmount:  so.TotalAmount,
			Status:       so.Status,
			CreatedBy:    so.CreatedBy,
			CreatedAt:    so.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    so.UpdatedAt.Format(time.RFC3339),
		})
	}

	return responses, totalCount, nil
}

func (s *SaleOrderService) GetSaleOrderByID(ctx context.Context, id uuid.UUID) (*dto.SaleOrderResponse, error) {
	saleOrder, err := s.saleOrderRepo.GetSaleOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.SaleOrderResponse{
		ID:           saleOrder.ID,
		OrderNumber:  saleOrder.OrderNumber,
		CustomerName: saleOrder.CustomerName,
		TotalAmount:  saleOrder.TotalAmount,
		Status:       saleOrder.Status,
		CreatedBy:    saleOrder.CreatedBy,
		CreatedAt:    saleOrder.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    saleOrder.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *SaleOrderService) UpdateSaleOrder(ctx context.Context, id uuid.UUID, req *dto.UpdateSaleOrderRequest) error {
	existingSaleOrder, err := s.saleOrderRepo.GetSaleOrderByID(ctx, id)
	if err != nil {
		return err
	}

	existingSaleOrder.CustomerName = req.CustomerName
	existingSaleOrder.TotalAmount = req.TotalAmount
	existingSaleOrder.Status = req.Status
	existingSaleOrder.UpdatedAt = time.Now()

	return s.saleOrderRepo.UpdateSaleOrder(ctx, existingSaleOrder)
}

func (s *SaleOrderService) DeleteSaleOrder(ctx context.Context, id uuid.UUID) error {
	return s.saleOrderRepo.DeleteSaleOrder(ctx, id)
}

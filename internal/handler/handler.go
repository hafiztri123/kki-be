package handler

import "github.com/hafiztri123/kki-be/internal/service"

type Handlers struct {
	UserHandler      *UserHandler
	SaleOrderHandler *SaleOrderHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		UserHandler:      NewUserHandler(services.UserService),
		SaleOrderHandler: NewSaleOrderHandler(services.SaleOrderService),
	}
}

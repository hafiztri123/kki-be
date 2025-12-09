package service

import "github.com/hafiztri123/kki-be/internal/repository"

type Services struct {
	UserService      *UserService
	SaleOrderService *SaleOrderService
}

func NewServices(repositories *repository.Repositories) *Services {
	return &Services{
		UserService:      NewUserService(repositories.UserRepository),
		SaleOrderService: NewSaleOrderService(repositories.SaleOrderRepository),
	}
}

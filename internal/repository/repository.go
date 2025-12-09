package repository

import "github.com/jackc/pgx/v5/pgxpool"



type Repositories struct {
	UserRepository      *UserRepository
	SaleOrderRepository *SaleOrderRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		UserRepository:      NewUserRepository(db),
		SaleOrderRepository: NewSaleOrderRepository(db),
	}
}

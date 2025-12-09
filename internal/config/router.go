package config

import (
	"net/http"

	"github.com/hafiztri123/kki-be/internal/constants"
	"github.com/hafiztri123/kki-be/internal/handler"
	"github.com/hafiztri123/kki-be/internal/middleware"
)

func NewRouter(handlers *handler.Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", handlers.UserHandler.RegisterHandler)
	mux.HandleFunc("POST /api/v1/auth/login", handlers.UserHandler.LoginHandler)
	mux.HandleFunc("POST /api/v1/auth/logout", middleware.JWTMiddleware(handlers.UserHandler.LogoutHandler))

	mux.HandleFunc("GET /api/v1/sale-orders",
		middleware.JWTMiddleware(
			middleware.RBACMiddleware(handlers.SaleOrderHandler.GetSaleOrdersHandler, constants.RoleCashier, constants.RoleOwner)))

	mux.HandleFunc("GET /api/v1/sale-orders/{id}",
		middleware.JWTMiddleware(
			middleware.RBACMiddleware(handlers.SaleOrderHandler.GetSaleOrderByIDHandler, constants.RoleCashier, constants.RoleOwner)))

	mux.HandleFunc("POST /api/v1/sale-orders",
		middleware.JWTMiddleware(
			middleware.RBACMiddleware(handlers.SaleOrderHandler.CreateSaleOrderHandler, constants.RoleCashier, constants.RoleOwner)))

	mux.HandleFunc("PUT /api/v1/sale-orders/{id}",
		middleware.JWTMiddleware(
			middleware.RBACMiddleware(handlers.SaleOrderHandler.UpdateSaleOrderHandler, constants.RoleCashier, constants.RoleOwner)))

	mux.HandleFunc("DELETE /api/v1/sale-orders/{id}",
		middleware.JWTMiddleware(
			middleware.RBACMiddleware(handlers.SaleOrderHandler.DeleteSaleOrderHandler, constants.RoleCashier, constants.RoleOwner)))

	// Cashier
	mux.HandleFunc("GET /api/v1/users/cashier",
		middleware.JWTMiddleware(
			middleware.RBACMiddleware(handlers.UserHandler.GetCashiersHandler, constants.RoleOwner)))

	mux.HandleFunc("GET /api/v1/users/cashier/{id}",
		middleware.JWTMiddleware(
			middleware.RBACMiddleware(handlers.UserHandler.GetCashierByIDHandler, constants.RoleOwner)))

	mux.HandleFunc("POST /api/v1/users/cashier",
		middleware.JWTMiddleware(
			middleware.RBACMiddleware(handlers.UserHandler.CreateCashierHandler, constants.RoleOwner)))

	mux.HandleFunc("PUT /api/v1/users/cashier/{id}",
		middleware.JWTMiddleware(
			middleware.RBACMiddleware(handlers.UserHandler.UpdateCashierHandler, constants.RoleOwner)))

	mux.HandleFunc("DELETE /api/v1/users/cashier/{id}",
		middleware.JWTMiddleware(
			middleware.RBACMiddleware(handlers.UserHandler.DeleteCashierHandler, constants.RoleOwner)))

	return mux
}
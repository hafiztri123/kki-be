package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	apperror "github.com/hafiztri123/kki-be/internal/app_error"
	"github.com/hafiztri123/kki-be/internal/constants"
	"github.com/hafiztri123/kki-be/internal/dto"
	"github.com/hafiztri123/kki-be/internal/service"
	"github.com/hafiztri123/kki-be/internal/utils"
)

type SaleOrderHandler struct {
	saleOrderService *service.SaleOrderService
}

func NewSaleOrderHandler(saleOrderService *service.SaleOrderService) *SaleOrderHandler {
	return &SaleOrderHandler{
		saleOrderService: saleOrderService,
	}
}

func (h *SaleOrderHandler) CreateSaleOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSaleOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.NewSlogFailToDecode(r, err)
		utils.NewJSONResponse(w, http.StatusBadRequest, constants.MsgStatusError, constants.MsgBadRequest, nil)
		return
	}

	createdBy, ok := r.Context().Value(constants.ClaimsKeyID).(uuid.UUID)
	if !ok {
		slog.ErrorContext(r.Context(), constants.MsgFailedToExtractValueFromJWT, "path", r.URL.Path )
		utils.NewJSONResponse(w, http.StatusUnauthorized, constants.MsgStatusError, constants.MsgUnauthorized, nil)
		return
	}

	err := h.saleOrderService.CreateSaleOrder(r.Context(), &req, createdBy)
	if err != nil {
		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(w, http.StatusInternalServerError, constants.MsgStatusError, constants.MsgInternalServerError, nil)
		return
	}

	utils.NewJSONResponse(w, http.StatusCreated, constants.MsgStatusSuccess, constants.MsgSuccessCreate, nil)
}

func (h *SaleOrderHandler) GetSaleOrdersHandler(w http.ResponseWriter, r *http.Request) {
	pagination := utils.ParsePaginationParams(r)
	offset := utils.CalculateOffset(pagination.Page, pagination.Limit)

	saleOrders, totalCount, err := h.saleOrderService.GetSaleOrders(r.Context(), pagination.Limit, offset)
	if err != nil {
		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(w, http.StatusInternalServerError, constants.MsgStatusError, constants.MsgInternalServerError, nil)
		return
	}

	paginatedResponse := utils.NewPaginatedResponse(saleOrders, totalCount, pagination.Page, pagination.Limit)

	utils.NewJSONResponse(w, http.StatusOK, constants.MsgStatusSuccess, constants.MsgSuccessRetrieve, paginatedResponse)
}

func (h *SaleOrderHandler) GetSaleOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewJSONResponse(w, http.StatusBadRequest, constants.MsgStatusError, constants.MsgBadRequest, nil)
		return
	}

	saleOrder, err := h.saleOrderService.GetSaleOrderByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			utils.NewJSONResponse(w, http.StatusNotFound, constants.MsgStatusError, constants.MsgNotFound, nil)
			return
		}

		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(w, http.StatusInternalServerError, constants.MsgStatusError, constants.MsgInternalServerError, nil)
		return
	}

	utils.NewJSONResponse(w, http.StatusOK, constants.MsgStatusSuccess, constants.MsgSuccessRetrieve, saleOrder)
}

func (h *SaleOrderHandler) UpdateSaleOrderHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewJSONResponse(w, http.StatusBadRequest, constants.MsgStatusError, constants.MsgBadRequest, nil)
		return
	}

	var req dto.UpdateSaleOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.NewSlogFailToDecode(r, err)
		utils.NewJSONResponse(w, http.StatusBadRequest, constants.MsgStatusError, constants.MsgBadRequest, nil)
		return
	}

	err = h.saleOrderService.UpdateSaleOrder(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			utils.NewJSONResponse(w, http.StatusNotFound, constants.MsgStatusError, constants.MsgNotFound, nil)
			return
		}

		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(w, http.StatusInternalServerError, constants.MsgStatusError, constants.MsgInternalServerError, nil)
		return
	}

	utils.NewJSONResponse(w, http.StatusOK, constants.MsgStatusSuccess, constants.MsgSuccessUpdate, nil)
}

func (h *SaleOrderHandler) DeleteSaleOrderHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewJSONResponse(w, http.StatusBadRequest, constants.MsgStatusError, constants.MsgBadRequest, nil)
		return
	}

	err = h.saleOrderService.DeleteSaleOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			utils.NewJSONResponse(w, http.StatusNotFound, constants.MsgStatusError, constants.MsgNotFound, nil)
			return
		}

		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(w, http.StatusInternalServerError, constants.MsgStatusError, constants.MsgInternalServerError, nil)
		return
	}

	utils.NewJSONResponse(w, http.StatusOK, constants.MsgStatusSuccess, constants.MsgSuccessDelete, nil)
}

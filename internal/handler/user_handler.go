package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	apperror "github.com/hafiztri123/kki-be/internal/app_error"
	"github.com/hafiztri123/kki-be/internal/constants"
	"github.com/hafiztri123/kki-be/internal/dto"
	"github.com/hafiztri123/kki-be/internal/service"
	"github.com/hafiztri123/kki-be/internal/utils"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.NewSlogFailToDecode(r, err)
		utils.NewJSONResponse(w, http.StatusBadRequest, constants.MsgBadRequest, constants.MsgBadRequest, nil)
		return
	}

	err := u.userService.Register(r.Context(), &req)
	if err != nil {
		if errors.Is(err, apperror.ErrEmailAlreadyExists) {
			utils.NewJSONResponse(
				w,
				http.StatusConflict,
				constants.MsgStatusError,
				constants.MsgEmailAlreadyExists,
				nil,
			)
			return
		}

		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(
			w,
			http.StatusInternalServerError,
			constants.MsgStatusError,
			constants.MsgInternalServerError,
			nil,
		)

		return
	}

	utils.NewJSONResponse(
		w,
		http.StatusCreated,
		constants.MsgStatusSuccess,
		constants.MsgSuccessRegister,
		nil,
	)
}

func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.NewSlogFailToDecode(r, err)
		utils.NewJSONResponse(w, http.StatusBadRequest, constants.MsgBadRequest, constants.MsgBadRequest, nil)
		return
	}

	res, err := u.userService.Login(r.Context(), &req)
	if err != nil {
		if errors.Is(err, apperror.ErrInvalidCredentials) || errors.Is(err, apperror.ErrNotFound) {
			utils.NewJSONResponse(
				w,
				http.StatusBadRequest,
				constants.MsgStatusError,
				constants.MsgInvalidCredentials,
				nil,
			)
			return

		}


		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(
			w,
			http.StatusInternalServerError,
			constants.MsgStatusError,
			constants.MsgInternalServerError,
			nil,
		)

		return
	}

	utils.NewJSONResponse(
		w,
		http.StatusOK,
		constants.MsgStatusSuccess,
		constants.MsgSuccessLogin,
		res,
	)
}


func (u *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// used when using refresh token or cookie
	utils.NewJSONResponse(
		w,
		http.StatusOK,
		constants.MsgStatusSuccess,
		constants.MsgSuccessLogout,
		nil,
	)
}

func (u *UserHandler) CreateCashierHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCashierRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.NewSlogFailToDecode(r, err)
		utils.NewJSONResponse(w, http.StatusBadRequest, constants.MsgStatusError, constants.MsgBadRequest, nil)
		return
	}

	err := u.userService.CreateCashier(r.Context(), &req)
	if err != nil {
		if errors.Is(err, apperror.ErrEmailAlreadyExists) {
			utils.NewJSONResponse(
				w,
				http.StatusConflict,
				constants.MsgStatusError,
				constants.MsgEmailAlreadyExists,
				nil,
			)
			return
		}

		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(
			w,
			http.StatusInternalServerError,
			constants.MsgStatusError,
			constants.MsgInternalServerError,
			nil,
		)
		return
	}

	utils.NewJSONResponse(
		w,
		http.StatusCreated,
		constants.MsgStatusSuccess,
		constants.MsgSuccessCreate,
		nil,
	)
}

func (u *UserHandler) GetCashiersHandler(w http.ResponseWriter, r *http.Request) {
	pagination := utils.ParsePaginationParams(r)
	offset := utils.CalculateOffset(pagination.Page, pagination.Limit)

	cashiers, totalCount, err := u.userService.GetCashiers(r.Context(), pagination.Limit, offset)
	if err != nil {
		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(w, http.StatusInternalServerError, constants.MsgStatusError, constants.MsgInternalServerError, nil)
		return
	}

	paginatedResponse := utils.NewPaginatedResponse(cashiers, totalCount, pagination.Page, pagination.Limit)

	utils.NewJSONResponse(w, http.StatusOK, constants.MsgStatusSuccess, constants.MsgSuccessRetrieve, paginatedResponse)
}

func (u *UserHandler) GetCashierByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	cashier, err := u.userService.GetCashierByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			utils.NewJSONResponse(w, http.StatusNotFound, constants.MsgStatusError, constants.MsgNotFound, nil)
			return
		}

		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(w, http.StatusInternalServerError, constants.MsgStatusError, constants.MsgInternalServerError, nil)
		return
	}

	utils.NewJSONResponse(w, http.StatusOK, constants.MsgStatusSuccess, constants.MsgSuccessRetrieve, cashier)
}

func (u *UserHandler) UpdateCashierHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req dto.UpdateCashierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.NewSlogFailToDecode(r, err)
		utils.NewJSONResponse(w, http.StatusBadRequest, constants.MsgStatusError, constants.MsgBadRequest, nil)
		return
	}

	err := u.userService.UpdateCashier(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			utils.NewJSONResponse(w, http.StatusNotFound, constants.MsgStatusError, constants.MsgNotFound, nil)
			return
		}

		if errors.Is(err, apperror.ErrEmailAlreadyExists) {
			utils.NewJSONResponse(
				w,
				http.StatusConflict,
				constants.MsgStatusError,
				constants.MsgEmailAlreadyExists,
				nil,
			)
			return
		}

		utils.NewSlogInternalServerError(r, err)
		utils.NewJSONResponse(w, http.StatusInternalServerError, constants.MsgStatusError, constants.MsgInternalServerError, nil)
		return
	}

	utils.NewJSONResponse(w, http.StatusOK, constants.MsgStatusSuccess, constants.MsgSuccessUpdate, nil)
}

func (u *UserHandler) DeleteCashierHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := u.userService.DeleteCashier(r.Context(), id)
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

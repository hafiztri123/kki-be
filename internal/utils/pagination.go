package utils

import (
	"math"
	"net/http"
	"strconv"

	"github.com/hafiztri123/kki-be/internal/dto"
)

func ParsePaginationParams(r *http.Request) dto.PaginationRequest {
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	return dto.PaginationRequest{
		Limit: limit,
		Page:  page,
	}
}

func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

func CalculateTotalPages(totalItems int64, limit int) int {
	return int(math.Ceil(float64(totalItems) / float64(limit)))
}

func NewPaginatedResponse(data interface{}, totalItems int64, page, limit int) dto.PaginatedResponse {
	return dto.PaginatedResponse{
		Data:       data,
		TotalItems: totalItems,
		TotalPages: CalculateTotalPages(totalItems, limit),
		Page:       page,
		Limit:      limit,
	}
}

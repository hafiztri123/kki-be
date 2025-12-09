package utils

import (
	"log/slog"
	"net/http"

	"github.com/hafiztri123/kki-be/internal/constants"
)



func NewSlogInternalServerError(r *http.Request, err error) {
	slog.ErrorContext(
		r.Context(),
		constants.MsgInternalServerError,
		"error", err.Error(),
		"path", r.URL.Path,
	)
}

func NewSlogFailToDecode(r *http.Request, err error) {
	slog.ErrorContext(
		r.Context(),
		constants.MsgFailDecode,
		"error", err.Error(),
		"path", r.URL.Path,
	)
}
package middleware

import (
	"net/http"
	"slices"

	"github.com/hafiztri123/kki-be/internal/constants"
	"github.com/hafiztri123/kki-be/internal/utils"
)

func RBACMiddleware(next http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(roles) == 0 {
			next(w, r)
		}

		userRole, ok := r.Context().Value(constants.ClaimsKeyRole).(string)

		if !ok || !slices.Contains(roles, userRole) {
			utils.NewJSONResponse(w, http.StatusForbidden, constants.MsgStatusError, constants.MsgForbidden, nil)
			return
		}

		next(w, r)
	}
}

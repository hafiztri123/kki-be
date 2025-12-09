package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hafiztri123/kki-be/internal/constants"
	"github.com/hafiztri123/kki-be/internal/utils"
)

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if !strings.HasPrefix(bearer, "Bearer ") || len(strings.Split(bearer, " ")) != 2 {
			fmt.Println("failed to parse header")

			utils.NewJSONResponse(w, http.StatusUnauthorized, constants.MsgStatusError, constants.MsgUnauthorized, nil)
			return
		}

		token := strings.Split(bearer, " ")[1]

		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.NewJSONResponse(w, http.StatusUnauthorized, constants.MsgStatusError, constants.MsgUnauthorized, nil)
			return
		}

		ctx := context.WithValue(r.Context(), constants.ClaimsKeyID, claims.Id)
		ctx = context.WithValue(ctx, constants.ClaimsKeyRole, claims.Role)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

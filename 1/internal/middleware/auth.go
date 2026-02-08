package middleware

import (
	"encoding/json"
	"net/http"

	"1/internal/models"
)

const apiKey = "secret123"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-KEY")

		if key != apiKey {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: "unauthorized",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

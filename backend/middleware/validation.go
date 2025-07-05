package middleware

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var validate = validator.New()

func ValidateJSON[T any](next func(http.ResponseWriter, *http.Request, T)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data T

		// Limit request body size
		r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		next(w, r, data)
	}
}

package middleware

import (
	"net/http"
	"temperature-backend/handler/util"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			err := util.HttpStatus{Code: http.StatusUnauthorized, Msg: "Authorization token is empty or not present in request"}
			http.Error(w, err.ToString(), err.Code)
			return
		}
		next.ServeHTTP(w, r)
	})
}

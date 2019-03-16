package middleware

import (
	"net/http"
	"temperature-backend/util"
)

func Post(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			err := util.HttpStatus{Code: http.StatusNotImplemented, Msg: "Not implemented"}
			http.Error(w, err.ToString(), err.Code)
			return
		}
		next.ServeHTTP(w, r)
	})
}
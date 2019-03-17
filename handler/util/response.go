package util

import (
	"encoding/json"
	"net/http"
)

func SetResponse(w http.ResponseWriter, resp interface{}) {
	respBody, err := json.Marshal(resp)
	if err != nil {
		HandleInternalError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		HandleInternalError(w, err)
	}
}

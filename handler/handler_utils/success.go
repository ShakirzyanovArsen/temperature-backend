package handler_utils

import (
	"encoding/json"
	"net/http"
	"temperature-backend/util"
)

func SetResponse(w http.ResponseWriter, resp interface{}) {
	respBody, err := json.Marshal(resp)
	if err != nil {
		util.HandleInternalError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		util.HandleInternalError(w, err)
	}
}

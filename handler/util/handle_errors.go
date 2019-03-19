package util

import (
	"log"
	"net/http"
	"temperature-backend/service"
)

func HandleInternalError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	httpErr := HttpStatus{Code: http.StatusInternalServerError, Msg: "Internal server error"}
	log.Print("Unexpected error occurred", err)
	http.Error(w, httpErr.ToString(), httpErr.Code)
}

func HandleResourceNotFound(w http.ResponseWriter) {
	http.Error(w, "404 not found", http.StatusNotFound)
}

func HandleServiceError(w http.ResponseWriter, e service.Error) {
	w.Header().Set("Content-Type", "application/json")
	if e.Code == http.StatusInternalServerError && e.Msg == "" {
		HandleInternalError(w, e)
	} else {
		status := HttpStatus{Code: e.Code, Msg: e.Msg}
		http.Error(w, status.ToString(), status.Code)
	}
}

func HandleSerializingError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	httpErr := HttpStatus{Code: http.StatusBadRequest, Msg: "Invalid json."}
	log.Print("Error occurred while unmarshaling.", err)
	http.Error(w, httpErr.ToString(), httpErr.Code)
}

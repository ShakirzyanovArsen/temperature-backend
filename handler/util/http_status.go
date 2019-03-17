package util

import (
	"encoding/json"
	"log"
)

type HttpStatus struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (e HttpStatus) ToString() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		log.Fatal("Error while http error marshaling")
	}
	return string(bytes)
}

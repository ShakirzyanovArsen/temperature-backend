package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func ParseJsonBody(w http.ResponseWriter, r *http.Request, entity interface{}) error {
	reqBuf := new(bytes.Buffer)
	_, err := reqBuf.ReadFrom(r.Body)
	if err != nil {
		HandleInternalError(w, err)
		return err
	}
	err = json.Unmarshal(reqBuf.Bytes(), &entity)
	if err != nil {
		HandleSerializingError(w, err)
		return err
	}
	return nil
}

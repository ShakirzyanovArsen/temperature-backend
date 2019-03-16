package functional_tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"temperature-backend/server"
	"testing"
)

func checkFieldsExists(body []byte, fields []string) (string, error){
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(body, &objmap)
	if err != nil {
		return "", err
	}
	for _, f := range fields {
		if objmap[f] == nil {
			return f, nil
		}
	}
	return "", nil
}

func setupServer () {
	server.Setup()
	http.ListenAndServe(":8080", nil)
}

func TestRegisterUser(t *testing.T) {
	url := "http://localhost:8080/user/register"
	tests := []TestStruct{
		{
			Name:               "success",
			RequestBody:        `{"email": "test@test"}`,
			ExpectedStatusCode: http.StatusCreated,
			ExistsFields: []string{"token"},
		},
		{
			Name:               "mail already in use",
			RequestBody:        `{"email": "test@test"}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExistsFields: []string{"code", "message"},
		},
	}
	go setupServer()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			reader := strings.NewReader(tt.RequestBody)
			req, err := http.NewRequest(http.MethodPost, url, reader)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}
			body, _ := ioutil.ReadAll(res.Body)
			if res.StatusCode != tt.ExpectedStatusCode {
				t.Error(fmt.Sprintf("wrong status code. expected %d, actual: %d", tt.ExpectedStatusCode, res.StatusCode))
			}
			fieldName, err := checkFieldsExists(body, tt.ExistsFields)
			if err != nil {
				t.Error("Cannot unmarshal response", string(body))
			}
			if fieldName != "" {
				t.Error(fmt.Sprintf("Field does not exists: %s", fieldName))
			}
		})
	}
}

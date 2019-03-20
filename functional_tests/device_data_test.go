package functional_tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestPushData(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	url := pushDataUrl
	existingUserEmail := "test@test.ru"
	_, err := createUser(existingUserEmail)
	if err != nil {
		t.Errorf("error while user create: %s", err)
	}
	validToken, err := registerDevice(existingUserEmail, "device1")
	if err != nil {
		t.Errorf("error while device register: %s", err)
	}
	validTokenMap := map[string]string{}
	validTokenMap["Authorization"] = validToken
	invalidTokenMap := map[string]string{}
	invalidTokenMap["Authorization"] = "invalidToken"
	validDateTime := "2019-03-17T21:08:00+05:00"
	tests := []TestStruct{
		{
			Name:               "success",
			RequestBody:        fmt.Sprintf(`{"date_time": "%s","temperature": 33.8}`, validDateTime),
			Headers:            validTokenMap,
			ExpectedStatusCode: http.StatusCreated,
			ExistsFields:       []string{"saved"},
		},
		{
			Name:               "wrong auth token",
			RequestBody:        fmt.Sprintf(`{"date_time": "%s","temperature": 33.8}`, validDateTime),
			Headers:            invalidTokenMap,
			ExpectedStatusCode: http.StatusUnauthorized,
			ExistsFields:       []string{"code", "message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			reader := strings.NewReader(tt.RequestBody)
			req, err := http.NewRequest(http.MethodPost, url, reader)
			for header, val := range tt.Headers {
				req.Header.Set(header, val)
			}
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

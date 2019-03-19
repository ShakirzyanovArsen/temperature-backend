package functional_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"temperature-backend/test_utils"
	"temperature-backend/view"
	"testing"
	"time"
)

func TestRegisterDevice(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test due to short mode")
	}
	url := registerDeviceUrl
	existingUserEmail := "test@test.ru"
	tests := []TestStruct{
		{
			Name:               "success",
			RequestBody:        fmt.Sprintf(`{"user_email": "%s", "device_name": "device1"}`, existingUserEmail),
			ExpectedStatusCode: http.StatusCreated,
			ExistsFields:       []string{"token"},
		},
		{
			Name:               "user does not exists",
			RequestBody:        `{"user_email": "user_not@exist.test", "device_name": "device1"}`,
			ExpectedStatusCode: http.StatusNotFound,
			ExistsFields:       []string{"code", "message"},
		},
	}
	srv := setupServer()
	defer func() {
		err := srv.Shutdown(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}()
	_, err := createUser(existingUserEmail)
	if err != nil {
		t.Errorf("error while user create: %s", err)
	}
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

func TestGetDeviceList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test due to short mode")
	}
	url := getDeviceListUrl
	srv := setupServer()
	defer func() {
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}()
	existingUserEmail := "list_test@test.ru"
	validTokenMap := map[string]string{}
	validUserToken, err := createUser(existingUserEmail)
	validTokenMap["Authorization"] = validUserToken
	tests := []TestStruct{
		{
			Name:               "success",
			RequestBody:        "",
			Headers:            validTokenMap,
			ExpectedStatusCode: http.StatusOK,
		},
	}
	d1Token, err := registerDevice(existingUserEmail, "device1")
	_, err = registerDevice(existingUserEmail, "device2")
	err = pushData(d1Token, "2019-03-17T21:08:00+05:00", 30.)
	err = pushData(d1Token, "2019-03-18T21:08:00+05:00", 40.)
	if err != nil {
		t.Errorf("error while pushing data: %s", err)
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			reader := strings.NewReader(tt.RequestBody)
			req, err := http.NewRequest(http.MethodGet, url, reader)
			req.Header.Set("Authorization", tt.Headers["Authorization"])
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
			resultView := view.DeviceListView{}
			err = json.Unmarshal(body, &resultView)
			if err != nil {
				t.Error("cant unmarshal expected result or response body")
			}
			test_utils.AssertInt(t, 2, len(resultView.Devices))
			test_utils.AssertString(t, "2019-03-18T21:08:00+05:00", resultView.Devices[0].LastDataTime)
			test_utils.AssertString(t, "", resultView.Devices[1].LastDataTime)
		})
	}
}

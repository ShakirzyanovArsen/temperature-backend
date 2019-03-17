package functional_tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"temperature-backend/server"
)

const createUserUrl = "http://localhost:8080/user/register"
const registerDeviceUrl = "http://localhost:8080/device/register"
const pushDataUrl = "http://localhost:8080/device/data"

type registerDeviceResp struct {
	Token string `json:"token"`
}

func checkFieldsExists(body []byte, fields []string) (string, error) {
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

func createUser(email string) error {
	reader := strings.NewReader(fmt.Sprintf(`{"email": "%s"}`, email))
	req, err := http.NewRequest(http.MethodPost, createUserUrl, reader)
	_, err = http.DefaultClient.Do(req)
	return err
}

func registerDevice(userEmail string, deviceName string) (string, error) {
	reader := strings.NewReader(fmt.Sprintf(`{"user_email": "%s", "device_name": "%s"}`, userEmail, deviceName))
	req, err := http.NewRequest(http.MethodPost, registerDeviceUrl, reader)
	res, err := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	tokenResult := registerDeviceResp{}
	err = json.Unmarshal(body, &tokenResult)
	return tokenResult.Token, err
}

func setupServer() {
	mux := server.Setup()
	http.ListenAndServe(":8080", mux)
}

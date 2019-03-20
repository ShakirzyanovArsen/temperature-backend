package functional_tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"temperature-backend/server"
	"temperature-backend/view"
	"time"
)

const createUserUrl = "http://localhost:8080/user/register"
const registerDeviceUrl = "http://localhost:8080/device/register"
const pushDataUrl = "http://localhost:8080/device/data"
const getDeviceListUrl = "http://localhost:8080/device/list"
const getDataListUrl = "http://localhost:8080/device/{id}"

type tokenResponse struct {
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

func createUser(email string) (string, error) {
	reader := strings.NewReader(fmt.Sprintf(`{"email": "%s"}`, email))
	req, err := http.NewRequest(http.MethodPost, createUserUrl, reader)
	res, err := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	tokenResult := tokenResponse{}
	err = json.Unmarshal(body, &tokenResult)
	return tokenResult.Token, err
}

func registerDevice(userEmail string, deviceName string) (string, error) {
	reader := strings.NewReader(fmt.Sprintf(`{"user_email": "%s", "device_name": "%s"}`, userEmail, deviceName))
	req, err := http.NewRequest(http.MethodPost, registerDeviceUrl, reader)
	res, err := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	tokenResult := tokenResponse{}
	err = json.Unmarshal(body, &tokenResult)
	return tokenResult.Token, err
}

func pushData(deviceToken string, dateTime string, temperature float64) error {
	reader := strings.NewReader(fmt.Sprintf(`{"date_time": "%s", "temperature": %f}`, dateTime, temperature))
	req, err := http.NewRequest(http.MethodPost, pushDataUrl, reader)
	req.Header.Set("Authorization", deviceToken)
	_, err = http.DefaultClient.Do(req)
	return err
}

func getDeviceList(userToken string) (view.DeviceListView, error) {
	reader := strings.NewReader("")
	req, err := http.NewRequest(http.MethodGet, getDeviceListUrl, reader)
	req.Header.Set("Authorization", userToken)
	res, err := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	listView := view.DeviceListView{}
	err = json.Unmarshal(body, &listView)
	return listView, err
}

func setupServer() http.Server {
	mux := server.Setup()
	srv := http.Server{Addr: ":8080", Handler: mux}
	time.Sleep(time.Second * 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(fmt.Sprintf("ListenAndServer(8080) falling: %s", err))
		}
	}()
	return srv
}

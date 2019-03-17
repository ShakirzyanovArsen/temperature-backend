package functional_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"temperature-backend/server"
)

const createUserUrl = "http://localhost:8080/user/register"
const createDeviceUrl = "http://localhost:8080/device/register"

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

func setupServer() {
	server.Setup()
	http.ListenAndServe(":8080", nil)
}

package api_test

import (
	models_v1 "app/api/models/v1"
	"app/api/status"
	"bytes"
	"encoding/json"
	"github.com/test-go/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestRegisterClient(t *testing.T) {
	cases := []struct {
		Name       string
		Method     string
		URL        string
		Body       interface{}
		Headers    map[string]string
		QueryParam map[string]string
		Response   interface{}

		ExpectedStatusCode int
		ExpectedBody       interface{}
	}{
		{
			Name:   "Bad Request",
			Method: http.MethodPost,
			URL:    "/api/auth/register",
			Body: models_v1.RegisterRequest{
				Name:     "Test",
				Phone:    "998956523212",
				Password: "test",
			},
			Response:           &models_v1.Response{},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody: models_v1.Response{
				Code:    status.StatusBadRequest.Code,
				Message: status.StatusBadRequest.Message,
			},
		},
		{
			Name:   "Invalid email",
			Method: http.MethodPost,
			URL:    "/api/auth/register",
			Body: models_v1.RegisterRequest{
				Name:     "Test",
				Email:    "me.com",
				Phone:    "998956523212",
				Password: "test",
			},
			Response:           &models_v1.Response{},
			ExpectedStatusCode: status.StatusBadEmail.Code,
			ExpectedBody: models_v1.Response{
				Code:    status.StatusBadEmail.Code,
				Message: status.StatusBadEmail.Message,
			},
		},
		{
			Name:   "Invalid phone",
			Method: http.MethodPost,
			URL:    "/api/auth/register",
			Body: models_v1.RegisterRequest{
				Name:     "Test",
				Email:    "me@g.com",
				Phone:    "997986543210",
				Password: "test",
			},
			Response:           &models_v1.Response{},
			ExpectedStatusCode: status.StatusBadPhone.Code,
			ExpectedBody: models_v1.Response{
				Code:    status.StatusBadPhone.Code,
				Message: status.StatusBadPhone.Message,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			res, err := MakeRequest(c.Method, c.URL, c.Body, c.Response, c.Headers, c.QueryParam)
			if assert.NoError(t, err) == false {
				t.Fatalf("got error while sending request: %v", err)
			}
			if assert.NotNil(t, res) == false {
				t.Fatalf("got nil response")
			}
			defer res.Body.Close()

			assert.Equal(t, c.ExpectedStatusCode, res.StatusCode)

			bodyBytes, err := ioutil.ReadAll(res.Body)
			if assert.NoError(t, err) == false {
				t.Fatalf("got error while ioutil.ReadAll: %v", err)
			}
			bodyString := string(bodyBytes)

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(c.ExpectedBody); err != nil {
				log.Fatalf("failed to encode struct: %v", err)
			}
			assert.Equal(t, strings.TrimSpace(buf.String()), bodyString)
		})
	}
}

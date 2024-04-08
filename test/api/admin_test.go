package api

import (
	models_v1 "app/api/models/v1"
	"github.com/bxcodec/faker/v3"
	"github.com/test-go/testify/assert"
	"net/http"
	"sync"
	"testing"
)

func TestAdmin(t *testing.T) {
	test_cases := []struct {
		Name               string
		Body               interface{}
		ExpectedStatusCode int
		ExpectedBody       interface{}
	}{
		{
			Name: "Bad request (no password)",
			Body: models_v1.RegisterRequest{
				Name:  "Test-admin",
				Email: faker.Email(),
				Phone: "998563201452",
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody: models_v1.Response{
				Code:    http.StatusBadRequest,
				Message: "Bad request",
			},
		},
	}

	wg := &sync.WaitGroup{}
	count := 0
	for _, c := range test_cases {
		wg.Add(1)

		go func() {
			defer wg.Done()
			t.Run(c.Name, createAdmin(t, c.Body, c.ExpectedBody, c.ExpectedStatusCode))
		}()
		count++
	}
	wg.Wait()

	t.Logf("Tested: %d", count)
}

func createAdmin(t *testing.T, req, res interface{}, code int) func(t *testing.T) {
	return func(t *testing.T) {
		respBody := models_v1.Response{}
		httpResponse, err := MakeRequest(http.MethodPost, "/api/admin", req, &respBody, nil, nil)
		if !assert.NoError(t, err) {
			t.Fatalf("got error while sending request: %v", err)
		}
		if !assert.NotNil(t, httpResponse) {
			t.Fatalf("got nil response")
		}
		defer httpResponse.Body.Close()

		assert.Equal(t, code, httpResponse.StatusCode)

		assert.Equal(t, respBody.Message, res.(models_v1.Response).Message)
	}
}

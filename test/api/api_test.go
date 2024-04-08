package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseUrl = "http://localhost:8123"
)

func MakeRequest(method, requestURL string, req, res interface{}, headers map[string]string, queryParams map[string]string) (*http.Response, error) {
	var bodyReader io.Reader

	switch b := req.(type) {
	case string:
		bodyReader = strings.NewReader(b)
	case []byte:
		bodyReader = bytes.NewReader(b)
	case url.Values: // For form data
		bodyReader = strings.NewReader(b.Encode())
		if headers != nil {
			headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
	case nil:
		// Nobody
	default: // Assume JSON
		jsonData, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonData)
		if headers != nil {
			headers["Content-Type"] = "application/json"
		}
	}

	request, err := http.NewRequest(method, baseUrl+requestURL, bodyReader)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	// Add query parameters
	if len(queryParams) > 0 {
		q := request.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		request.URL.RawQuery = q.Encode()
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp_body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(resp_body, res); err != nil {
		return resp, err
	}
	return resp, nil
}

package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// NewMicroserviceJSONRequest returns a response for the given request parameters.
// DONT FORGET TO CALL response.Body.Close() !!!!
//
// i.e. 'defer response.Body.Close()'
//
// token - an authorization header token to send with request
//
// request - the json body name/value pairs as map[string]any
//
// url - request address
//
// method - GET, POST, PUT, DELETE, etc
func NewMicroserviceJSONRequest(token string, requestBody map[string]any, url string, method string) (code int, response interface{}, err error) {

	//create a JSON object to include in the request body
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return http.StatusInternalServerError, response, err
	}
	//create a new buffer and write the JSON bytes to it
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return http.StatusInternalServerError, response, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-type", "application/json")
	}

	//make the request

	//create an HTTP client and send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, response, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return http.StatusInternalServerError, response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return http.StatusInternalServerError, response, err
	}

	return http.StatusOK, response, err
}

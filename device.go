package spark

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Device struct {
	ID    string
	Token AuthToken
}

type DeviceConfiguration func(*Device) error

// Create a new core instance
func NewDevice(id string, configs ...DeviceConfiguration) (*Device, error) {

	d := &Device{}

	d.ID = id

	for _, config := range configs {
		if err := config(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

///////////////////////////////////////////////////////////////
// Configuration functions
///////////////////////////////////////////////////////////////

// Add an auth token to a device object
func DeviceAuthToken(token AuthToken) DeviceConfiguration {

	return func(d *Device) error {
		d.Token = token
		return nil
	}
}

///////////////////////////////////////////////////////////////
// API methods
///////////////////////////////////////////////////////////////

// Get a variable from the Spark Cloud
func (c *Device) GetWithToken(name string, auth_token AuthToken) (*VariableResponse, error) {

	// Try and generate a token
	token, err := auth_token.Token()

	if err != nil {
		return nil, err
	}

	// Get the final request URL
	url := c.requestURL(name, token)

	// Create a client request
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	// Make sure the content type is correct
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Execute the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	// Ensure the body is closed after the request
	defer resp.Body.Close()

	// A response-or-error container
	var response struct {
		VariableResponse
		ErrorResponse
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return nil, err
	}

	if response.Err != "" {
		return nil, response.ErrorResponse
	}

	return &response.VariableResponse, nil
}

func (c *Device) CallWithToken(name string, auth_token AuthToken, args ...interface{}) (*FunctionResponse, error) {

	// Try and generate a token
	token, err := auth_token.Token()

	if err != nil {
		return nil, err
	}

	// Get the final request URL
	uri := c.requestURL(name, token)

	// A form values container
	form := url.Values{}

	// Collection of arguments
	a := make([]string, 0)

	for _, arg := range args {
		a = append(a, fmt.Sprintf("%v", arg))
	}

	// Send all arguments separated by a comma
	form.Set("args", strings.Join(a, ","))

	req, err := http.NewRequest("POST", uri, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response struct {
		FunctionResponse
		ErrorResponse
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return nil, err
	}

	if response.Err != "" {
		return nil, response.ErrorResponse
	}

	return &response.FunctionResponse, nil

}

func (c *Device) requestURL(terminus, token string) string {

	return Endpoint(&APIUrl{
		BaseUrl,
		APIVersion,
		"/devices/" + c.ID + "/" + terminus + "?access_token=" + token,
	})
}

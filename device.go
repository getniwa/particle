package spark

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Device struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Connected    bool              `json:"connected"`
	Variables    map[string]string `json:"variables"`
	Functions    []string          `json:"functions"`
	PatchVersion string            `json:"cc3000_patch_version"`

	AuthTokenProvider AuthTokenProvider `json:"-"`
	Token             AuthToken         `json:"-"`
}

// A variadiac configuration function for device object instantiation
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

// Add an auth token provider (AccessTokenSerive or AuthToken) to
// a device object
//
func DeviceAuthTokenProvider(provider AuthTokenProvider) DeviceConfiguration {

	return func(d *Device) error {
		d.AuthTokenProvider = provider
		return nil
	}
}

///////////////////////////////////////////////////////////////
// Private methods
///////////////////////////////////////////////////////////////

func (d *Device) token() (token AuthToken, err error) {

	if d.AuthTokenProvider != nil {
		token, err = d.AuthTokenProvider.AuthToken()

		if err != nil {
			return
		}

	} else if d.Token != nil {
		token = d.Token
	}

	if token == nil {
		return nil, fmt.Errorf("No access token or provider set")
	}

	return
}

func (c *Device) requestURL(terminus, token string) string {

	url := &APIUrl{
		BaseUrl,
		APIVersion,
		"/devices/" + c.ID + "/" + terminus + "?access_token=" + token,
	}

	return url.String()
}

///////////////////////////////////////////////////////////////
// API methods
///////////////////////////////////////////////////////////////

// Get a variable from the Spark Cloud
func (d *Device) Get(name string) (*VariableResponse, error) {

	token, err := d.token()

	if err != nil {
		return nil, fmt.Errorf("[Device.Get] Token error: %s", err)
	}

	return d.GetWithToken(name, token)
}

// Get a variable from the Spark Cloud using an auth token
func (c *Device) GetWithToken(name string, token AuthToken) (*VariableResponse, error) {

	if err := token.Valid(); err != nil {
		return nil, err
	}

	// Get the final request URL
	url := c.requestURL(name, token.String())

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

// Get a variable from the Spark Cloud
func (d *Device) Call(name string, args ...interface{}) (*FunctionResponse, error) {

	token, err := d.token()

	if err != nil {
		return nil, fmt.Errorf("[Device.Call] Token error: %s", err)
	}

	return d.CallWithToken(name, token, args...)
}

// Call a function on the device using an auth token
func (c *Device) CallWithToken(name string, token AuthToken, args ...interface{}) (*FunctionResponse, error) {

	if err := token.Valid(); err != nil {
		return nil, err
	}

	// Get the final request URL
	uri := c.requestURL(name, token.String())

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

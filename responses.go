package spark

import "fmt"

// CoreInfo is a representation of the json struct coreInfo.
type CoreInfo struct {
	LastApp   string `json:"last_app,omitempty"`
	LastHeard string `json:"last_heard,omitempty"`
	Connected bool   `json:"connected,omitempty"`
	DeviceID  string `json:"deviceID,omitempty"`
}

// A representation of the response
// received on doing a GET on a Device Variable.
type VariableResponse struct {
	Cmd      string      `json:"cmd,omitempty"`
	Name     string      `json:"name,omitempty"`
	Result   interface{} `json:"result,omitempty"`
	CoreInfo CoreInfo    `json:"coreInfo,omitempty"`
}

// The representation of a response when a
// particular function is invoked via the REST API.
type FunctionResponse struct {
	DeviceID    string `json:"id"`
	Name        string `json:"name"`
	Connected   bool   `json:"connected"`
	ReturnValue int32  `json:"return_value"`
}

// An error response
type ErrorResponse struct {
	Code             int32  `json:"code, omitempty"`
	Err              string `json:"error, omitempty"`
	ErrorDescription string `json:"error_description, omitempty"`
	Info             string `json:"info, omitempty"`
}

func (e ErrorResponse) Error() string {
	errorMsg := ""

	// Just hoping that there isn't any error code of 0
	if e.Code != 0 {
		errorMsg = fmt.Sprintf(
			"Error Code: %v, Error: %v, ", e.Code, e.Err)
	}
	if e.ErrorDescription != "" {
		errorMsg += fmt.Sprintf("Error Description: %v",
			e.ErrorDescription)
	}
	if e.Info != "" {
		errorMsg += fmt.Sprint("Info: %v", e.Info)
	}

	return errorMsg
}

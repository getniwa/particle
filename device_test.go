package spark

import "testing"

const (
	ExampleVariable = "version"
	ExampleFunc     = "toggleLamp"
)

func validDevice(t *testing.T) *Device {

	d, creation_error := NewDevice(CoreID)

	if creation_error != nil {
		t.Fatalf("NewDevice: %s", creation_error)
	}

	return d
}

func validDeviceWithToken(t *testing.T, token AuthToken) *Device {

	d, creation_error := NewDevice(
		CoreID,
		DeviceAuthToken(token),
	)

	if creation_error != nil {
		t.Fatalf("NewDevice: %s", creation_error)
	}

	if d.Token != token {
		t.Fatalf("validDeviceWithToken: Tokens don't match")
	}

	return d
}

func validDeviceWithTokenProvider(t *testing.T, p AuthTokenProvider) *Device {

	d, creation_error := NewDevice(
		CoreID,
		DeviceAuthTokenProvider(p),
	)

	if creation_error != nil {
		t.Fatalf("NewDevice: %s", creation_error)
	}

	if d.AuthTokenProvider != p {
		t.Fatalf("validDeviceWithTokenProvider: Tokens don't match")
	}

	return d
}

func Test_GetVariableExplicitToken(t *testing.T) {

	core := validDevice(t)
	token := validAuthTokenFromService(t)

	if _, err := core.GetWithToken(ExampleVariable, token); err != nil {
		t.Error(err)
	}
}

func Test_CallFuncExplicitToken(t *testing.T) {

	core := validDevice(t)
	token := validAuthTokenFromService(t)

	if _, err := core.CallWithToken(ExampleFunc, token, 1, 2); err != nil {
		t.Error(err)
	}
}

func Test_GetVariableImplicitToken(t *testing.T) {

	token := validAuthTokenFromService(t)
	core := validDeviceWithToken(t, token)

	if _, err := core.Get(ExampleVariable); err != nil {
		t.Error(err)
	}
}

func Test_CallFuncImplicitToken(t *testing.T) {

	token := validAuthTokenFromService(t)
	core := validDeviceWithToken(t, token)

	if _, err := core.Call(ExampleFunc, 1, 2); err != nil {
		t.Error(err)
	}
}

func Test_GetVariableProvider(t *testing.T) {

	core := validDeviceWithTokenProvider(t, aTokenService)

	if _, err := core.Get(ExampleVariable); err != nil {
		t.Error(err)
	}
}

func Test_CallFuncProvider(t *testing.T) {

	core := validDeviceWithTokenProvider(t, aTokenService)

	if _, err := core.Call(ExampleFunc, 1, 2); err != nil {
		t.Error(err)
	}
}

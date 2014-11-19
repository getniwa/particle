package spark

import "testing"

const (
	ExampleVariable = "version"
	ExampleFunc     = "toggleLamp"
)

func Test_GetVariableExplicitToken(t *testing.T) {

	core, creation_error := NewDevice(CoreID)

	if creation_error != nil {
		t.Fatalf("NewDevice: %s", creation_error)
	}

	token, err := aTokenService.GetAccessToken()

	if err != nil {
		t.Error(err)
	}

	if _, err := core.GetWithToken(ExampleVariable, token); err != nil {
		t.Error(err)
	}
}

func Test_CallFuncExplicitToken(t *testing.T) {

	core, creation_error := NewDevice(CoreID)

	if creation_error != nil {
		t.Fatalf("NewDevice: %s", creation_error)
	}

	token, err := aTokenService.GetAccessToken()

	if err != nil {
		t.Error(err)
	}

	if _, err := core.CallWithToken(ExampleFunc, token, 1, 2); err != nil {
		t.Error(err)
	}
}

package spark

import "testing"

const (
	ExampleVariable = "version"
	ExampleFunc     = "toggleLamp"
)

func Test_GetVariable(t *testing.T) {

	core := NewCore(CoreID)

	token, err := aTokenService.GetAccessToken()

	if err != nil {
		t.Error(err)
	}

	if _, err := core.Get(ExampleVariable, token); err != nil {
		t.Error(err)
	}
}

func Test_CallFunc(t *testing.T) {

	core := NewCore(CoreID)

	token, err := aTokenService.GetAccessToken()

	if err != nil {
		t.Error(err)
	}

	if _, err := core.Call(ExampleFunc, token, 1, 2); err != nil {
		t.Error(err)
	}
}

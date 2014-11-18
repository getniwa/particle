package spark

import "testing"

func Test_GetVariable(t *testing.T) {

	core := NewCore(CoreID)

	token, err := aTokenService.GetAccessToken()

	if err != nil {
		t.Error(err)
	}

	//defer aTokenService.DeleteAccessToken(token)

	if _, err := core.Get("version", token); err != nil {
		t.Error(err)
	}
}

func Test_CallFunc(t *testing.T) {

	core := NewCore(CoreID)

	token, err := aTokenService.GetAccessToken()

	if err != nil {
		t.Error(err)
	}

	//defer aTokenService.DeleteAccessToken(token)

	if _, err := core.Call("toggleLamp", token, 1, 2); err != nil {
		t.Error(err)
	}
}

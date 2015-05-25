package particle

import "testing"

func validAuthTokenFromService(t *testing.T) AuthToken {

	token, err := aTokenService.AuthToken()

	if err != nil {
		t.Fatalf("GetAccessToken failed: %s", err)
	}

	return token
}

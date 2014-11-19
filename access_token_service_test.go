package spark

import "testing"

var aTokenService = NewAccessTokenService(UserName, Password)

func Test_GetAccessTokenHappyPath(t *testing.T) {

	if _, err := aTokenService.AuthToken(); err != nil {
		t.Fatalf("Can't get access token: %s", err)
	}
}

func Test_ListAllAccessTokens(t *testing.T) {

	_, err := aTokenService.ListAllAccessTokens()

	if err != nil {
		t.Error(err)
	}
}

// This is deprecated for the moment, because the DELETE
// method on the Spark Cloud API seems to take forever
func Test_DeleteAccessToken(t *testing.T) {
	/*
		delTokenResp, err := aTokenService.DeleteAccessToken(
			aTokenService.TokenList[2],
		)

		if err != nil {
			t.Error(err)
		}

		if !delTokenResp.Status {
			t.Error("Failed to delete token.")
		}
	*/
}

func Test_GetAccessTokenSadPath(t *testing.T) {

	badTokenService := NewAccessTokenService("Invalid-username", "Wrong-password")

	if _, err := badTokenService.AuthToken(); err == nil {
		t.Fatalf("Expected token service failure, got success")
	}
}

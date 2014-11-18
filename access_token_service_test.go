package spark

import "testing"

var aTokenService = NewAccessTokenService(UserName, Password)

func Test_GetAccessTokenHappyPath(t *testing.T) {

	if _, err := aTokenService.GetAccessToken(); err != nil {
		t.Fatalf("Can't get access token: %s", err)
	}
}

func Test_ListAllAccessTokens(t *testing.T) {

	if _, err := aTokenService.ListAllAccessTokens(); err != nil {
		t.Error(err)
	}
}

func Test_DeleteAccessToken(t *testing.T) {

	// Not touching 0th and 1st tokens in list of tokens as they
	// are for "user" and "spark-ide". Also, there will always be
	// a third token as TestListAllAccessTokens is called first.
	// TODO: Might change if using Mock.
	delTokenResp, err := aTokenService.DeleteAccessToken(
		aTokenService.TokenList.Tokens[2],
	)

	if err != nil {
		t.Error(err)
	}

	if !delTokenResp.Status {
		t.Error("Failed to delete token.")
	}
}

func Test_GetAccessTokenSadPath(t *testing.T) {

	badTokenService := NewAccessTokenService("Invalid-username", "Wrong-password")

	if _, err := badTokenService.GetAccessToken(); err == nil {
		t.Fatalf("Expected token service failure, got success")
	}
}

package spark

import "fmt"

type AccessToken struct {
	TokenValue string `json:"token"`
	ExpiresAt  string `json:"expires_at"`
	Client     string `json:"client"`
}

func NewAccessToken() *AccessToken {
	aToken := &AccessToken{}
	return aToken
}

type AccessTokenList struct {
	Tokens []*AccessToken
}

func NewAccessTokenList() *AccessTokenList {
	aTokenList := &AccessTokenList{}
	aTokenList.Tokens = []*AccessToken{}
	return aTokenList
}

func (t *AccessToken) Valid() error {

	if len(t.TokenValue) == 0 {
		return fmt.Errorf("Expected an access token of non-zero length")
	}

	return nil
}

func (t *AccessToken) Token() (string, error) {

	if err := t.Valid(); err != nil {
		return "", err
	}

	return t.TokenValue, nil
}

package spark

import "fmt"

const (
	GRANT_TYPE                 = "password"
	EXPECTED_ACCESS_TOKEN_TYPE = "bearer"
	DEFAULT_TOKEN_CLIENT       = "spark"
	EXPIRY_TIME_THRESHOLD      = 10
)

// OAuthRequest is a struct for the body of a request to get an access
// token.
type OAuthRequest struct {
	GrantType string `json:"grant_type, omitempty"`
	UserName  string `json:"username, omitoempty"`
	Password  string `json:"password, omitempty"`
}

func NewOAuthRequest(username, password string) *OAuthRequest {
	oaReq := &OAuthRequest{}
	oaReq.GrantType = GRANT_TYPE
	oaReq.UserName = username
	oaReq.Password = password

	return oaReq
}

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (t *OAuthResponse) Valid() error {

	if len(t.AccessToken) == 0 {
		return fmt.Errorf("Expected an access token of non-zero length")
	}

	if g, e := t.TokenType, EXPECTED_ACCESS_TOKEN_TYPE; g != e {
		return fmt.Errorf("Expected token type %s, got %s", e, g)
	}

	if e := t.ExpiresIn; e < EXPIRY_TIME_THRESHOLD {
		return fmt.Errorf("Auth token expiry time is too short (%d seconds)", e)
	}

	return nil
}

func (t *OAuthResponse) Token() (string, error) {

	if err := t.Valid(); err != nil {
		return "", err
	}

	return t.AccessToken, nil
}

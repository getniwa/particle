package spark

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AccessTokenService is an entrypoint for any AccessToken related operation.
type AccessTokenService struct {
	OaRequest *OAuthRequest
	Token     *AccessToken
}

// Expected response from a delete-access-token request
type DeleteAccessTokenResponse struct {
	Status bool `json:"ok"`
}

// Create a new access token service with a given username and password.
// This is passive, and won't return an error if the credentials are
// incorrect.
//
func NewAccessTokenService(username, password string) *AccessTokenService {

	aTokenService := &AccessTokenService{}
	aTokenService.OaRequest = NewOAuthRequest(username, password)

	return aTokenService
}

// Fetch the current access token, if it's available and not expired.
// If not, look in the list of existing tokens. If that fails, create a
// new one.
//
func (s *AccessTokenService) GetAccessToken() (AuthToken, error) {

	if s.Token != nil {
		// Check to see if it has expired
		if s.Token.ExpiresAt.After(time.Now()) {
			return s.Token, nil
		}
	}

	// Otherwise, check the list
	list, err := s.ListAllAccessTokens()

	if err != nil {
		return nil, fmt.Errorf("[GetAccessToken]: ListAllAccessTokens failed (%s)", err)
	}

	for _, token := range list {

		// Only use 'spark' token
		if token.Client != DEFAULT_TOKEN_CLIENT {
			continue
		}

		if token.ExpiresAt.After(time.Now()) {

			// Store the token
			s.Token = token

			// Return the pointer
			return s.Token, nil
		}
	}

	// Finally, if it's not anywhere, create a new one
	s.Token = &AccessToken{}
	response, err2 := s.CreateAccessToken()

	if err2 != nil {
		return nil, fmt.Errorf("[CreateAccessToken]: ListAllAccessTokens failed (%s)", err2)
	}

	s.Token.TokenValue = response.AccessToken
	s.Token.ExpiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	s.Token.Client = DEFAULT_TOKEN_CLIENT

	return s.Token, nil
}

// Returns an AccessToken in the form of a OAuthResponse object.
func (s *AccessTokenService) CreateAccessToken() (*OAuthResponse, error) {

	urlStr := Endpoint(
		&APIUrl{
			BaseUrl:  BaseUrl,
			Endpoint: "/oauth/token"},
	)

	form := url.Values{}
	form.Set("grant_type", s.OaRequest.GrantType)
	form.Set("username", s.OaRequest.UserName)
	form.Set("password", s.OaRequest.Password)

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	// Set some headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(BasicAuthId, BasicAuthPassword)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	oauthResp := &OAuthResponse{}

	if err = json.NewDecoder(resp.Body).Decode(oauthResp); err != nil {
		return nil, err
	}

	if err := oauthResp.Valid(); err != nil {
		return nil, fmt.Errorf("Token is not valid: %s", err)
	}

	return oauthResp, nil
}

func (s *AccessTokenService) ListAllAccessTokens() ([]*AccessToken, error) {

	urlStr := Endpoint(&APIUrl{BaseUrl, APIVersion, "/access_tokens"})

	req, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(s.OaRequest.UserName, s.OaRequest.Password)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Create an access token slice
	tokens := make([]*AccessToken, 0)

	if err = json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AccessTokenService) DeleteAccessToken(a AuthToken) (*DeleteAccessTokenResponse, error) {

	token, err := a.Token()

	if err != nil {
		return nil, err
	}

	ep := "/access_tokens/" + token
	urlStr := Endpoint(&APIUrl{BaseUrl, APIVersion, ep})

	req, err := http.NewRequest("DELETE", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(s.OaRequest.UserName, s.OaRequest.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	delTokenResp := &DeleteAccessTokenResponse{}

	if err := json.NewDecoder(resp.Body).Decode(delTokenResp); err != nil {
		return nil, err
	}

	return delTokenResp, nil
}

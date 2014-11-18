package spark

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// AccessTokenService is an entrypoint for any AccessToken related operation.
type AccessTokenService struct {
	OaRequest *OAuthRequest
	Token     *AccessToken
	TokenList *AccessTokenList
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
	aTokenService.Token = NewAccessToken()
	aTokenService.TokenList = NewAccessTokenList()

	return aTokenService
}

// Returns an AccessToken in the form of a OAuthResponse object.
func (s *AccessTokenService) GetAccessToken() (*OAuthResponse, error) {

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

func (s *AccessTokenService) ListAllAccessTokens() (*AccessTokenList, error) {

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

	s.TokenList = NewAccessTokenList()

	if err = json.NewDecoder(resp.Body).Decode(&s.TokenList.Tokens); err != nil {
		return nil, err
	}

	return s.TokenList, nil
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

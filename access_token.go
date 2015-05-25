package spark

import (
	"fmt"
	"time"
)

type AccessToken struct {
	TokenValue string     `json:"token"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Client     string     `json:"client"`
}

func NewAccessToken() *AccessToken {
	aToken := &AccessToken{}
	return aToken
}
func (t *AccessToken) Valid() error {

	if len(t.TokenValue) == 0 {
		return fmt.Errorf("Expected an access token of non-zero length")
	}

	if t.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("Token expired at %s", t.ExpiresAt.Format(time.RFC3339))
	}

	return nil
}

func (t *AccessToken) String() string {
	return t.TokenValue
}

func (t *AccessToken) AuthToken() (AuthToken, error) {

	if err := t.Valid(); err != nil {
		return nil, err
	}

	return t, nil
}

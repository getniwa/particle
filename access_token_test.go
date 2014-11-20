package spark

import (
	"testing"
	"time"
)

func Test_AccessTokenValidity(t *testing.T) {

	responses := []struct {
		response *AccessToken
		expected bool
	}{
		{
			&AccessToken{
				"SOMETOKENSSTRING",
				time.Now().Add(1 * time.Hour),
				DEFAULT_TOKEN_CLIENT,
			},
			true,
		},
		{
			&AccessToken{
				"",
				time.Now().Add(1 * time.Hour),
				DEFAULT_TOKEN_CLIENT,
			},
			false,
		},
		{
			&AccessToken{
				"SOMETOKENSSTRING",
				time.Now().Add(-1 * time.Hour),
				DEFAULT_TOKEN_CLIENT,
			},
			false,
		},
	}

	for i, r := range responses {

		err := r.response.Valid()
		outcome := true

		if err != nil {
			outcome = false
		}

		if outcome != r.expected {
			t.Errorf("Response [%d]: expected %v, got %v (err: %s)", i, r.expected, outcome, err)
		}
	}
}

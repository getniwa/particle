package spark

import (
	"testing"
	"time"
)

func Test_AccessTokenValidity(t *testing.T) {

	in_one_hour, one_hour_ago := time.Now(), time.Now()

	in_one_hour = in_one_hour.Add(1 * time.Hour)
	one_hour_ago = one_hour_ago.Add(-1 * time.Hour)

	responses := []struct {
		response *AccessToken
		expected bool
	}{
		{
			&AccessToken{
				"SOMETOKENSSTRING",
				&in_one_hour,
				DEFAULT_TOKEN_CLIENT,
			},
			true,
		},
		{
			&AccessToken{
				"",
				&in_one_hour,
				DEFAULT_TOKEN_CLIENT,
			},
			false,
		},
		{
			&AccessToken{
				"SOMETOKENSSTRING",
				&one_hour_ago,
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

package particle

import "testing"

func Test_OuthValidity(t *testing.T) {

	responses := []struct {
		response *OAuthResponse
		expected bool
	}{
		{
			&OAuthResponse{
				"SOMETOKENSSTRING",
				EXPECTED_ACCESS_TOKEN_TYPE,
				77000,
			},
			true,
		},
		{
			&OAuthResponse{
				"SOMETOKENSSTRING",
				EXPECTED_ACCESS_TOKEN_TYPE,
				0,
			},
			false,
		},
		{
			&OAuthResponse{
				"",
				EXPECTED_ACCESS_TOKEN_TYPE,
				1000,
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

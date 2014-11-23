package spark

import (
	"testing"
	"time"
)

func Test_EventDeviceLoad(t *testing.T) {

	e := Event{
		Name:      EVENT_NAME,
		Data:      EVENT_DATA,
		TTL:       "60",
		Published: time.Now(),
		CoreID:    CoreID,
	}

	token := validAuthTokenFromService(t)

	d, creation_error := e.Device(
		DeviceAuthToken(token),
	)

	if creation_error != nil {
		t.Fatalf("Event.Device: %s", creation_error)
	}

	if d.Token != token {
		t.Fatalf("validDeviceWithToken: Tokens don't match")
	}

}

package particle

import (
	"fmt"
	"time"
)

// Represents a Server-Sent Event
type Event struct {
	Name      string
	Data      string    `json:"data"`
	TTL       string    `json:"ttl"`
	Published time.Time `json:"published_at"`
	CoreID    string    `json:"coreid"`
}

func (e *Event) Device(dc ...DeviceConfiguration) (*Device, error) {

	if e.CoreID == "" {
		return nil, fmt.Errorf("CoreID must not be an empty string")
	}

	return NewDevice(e.CoreID, dc...)
}

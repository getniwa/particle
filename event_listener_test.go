package spark

import "testing"

const (
	EVENT_NAME = "test-successful"
	EVENT_DATA = "The test worked!"
)

func Test_Listener(t *testing.T) {

	l, err := NewEventListener(aTokenService)

	if err != nil {
		t.Fatalf("NewEventListener: %s", err)
	}

	// Run the listener
	go l.Listen()

	// Use the custom firmware endpoint to push something
	// back into the listener
	device := validDeviceWithTokenProvider(t, aTokenService)
	device.Call(PublishFunc)

	for event := range l.OutputChan {

		if event.Name == EVENT_NAME {

			if g, e := event.Data, EVENT_DATA; g != e {
				t.Fatalf("event.Data: got %s, expected %s", g, e)
			}

			break
		}
	}

	l.Stop()
}

func Test_ListenerForDevice(t *testing.T) {

	device := validDeviceWithTokenProvider(t, aTokenService)
	l, err := NewEventListener(aTokenService, EventListenerDevice(device))

	if err != nil {
		t.Fatalf("NewEventListener: %s", err)
	}

	// Run the listener
	go l.Listen()

	// Use the custom firmware endpoint to push something
	// back into the listener
	device.Call(PublishFunc)

	for event := range l.OutputChan {

		if event.Name == EVENT_NAME {

			if g, e := event.Data, EVENT_DATA; g != e {
				t.Fatalf("event.Data: got %s, expected %s", g, e)
			}

			break
		}

	}

	l.Stop()
}

func Test_ListenerForDeviceWithPrefix(t *testing.T) {

	device := validDeviceWithTokenProvider(t, aTokenService)
	l, err := NewEventListener(
		aTokenService,
		EventListenerDevice(device),
		EventListenerPrefix("test"),
	)

	if err != nil {
		t.Fatalf("NewEventListener: %s", err)
	}

	// Run the listener
	go l.Listen()

	// Use the custom firmware endpoint to push something
	// back into the listener
	device.Call(PublishFunc)

	for event := range l.OutputChan {

		if event.Name == EVENT_NAME {

			if g, e := event.Data, EVENT_DATA; g != e {
				t.Fatalf("event.Data: got %s, expected %s", g, e)
			}

			break
		}
	}

	l.Stop()
}

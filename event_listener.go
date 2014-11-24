package spark

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type EventChannel chan Event
type ErrorChannel chan error
type EventListenerConfiguration func(*EventListener) error

// A subscriber that listens to publication events
type EventListener struct {
	OutputChan        EventChannel
	ErrorChan         ErrorChannel
	AuthTokenProvider AuthTokenProvider
	Device            *Device
	prefix            string
	httpResponse      *http.Response
	running           bool
}

func NewEventListener(provider AuthTokenProvider, configs ...EventListenerConfiguration) (*EventListener, error) {

	e := &EventListener{
		AuthTokenProvider: provider,
	}

	for _, config := range configs {
		if err := config(e); err != nil {
			return nil, err
		}
	}

	if e.AuthTokenProvider == nil {
		return nil, fmt.Errorf("Auth token provider must not be nil")
	}

	if e.OutputChan == nil {
		e.OutputChan = make(chan Event)
	}

	if e.ErrorChan == nil {
		e.ErrorChan = make(chan error)
	}

	if e.httpResponse == nil {
		if err := e.connect(); err != nil {
			return nil, err
		}
	}

	return e, nil
}

///////////////////////////////////////////////////////////////////////////
// Config methods
///////////////////////////////////////////////////////////////////////////

func EventListenerOutputChannel(c EventChannel) EventListenerConfiguration {

	return func(e *EventListener) error {
		e.OutputChan = c
		return nil
	}
}

func EventListenerErrorChannel(c ErrorChannel) EventListenerConfiguration {

	return func(e *EventListener) error {
		e.ErrorChan = c
		return nil
	}
}

func EventListenerDevice(d *Device) EventListenerConfiguration {

	return func(e *EventListener) error {
		e.Device = d
		return nil
	}
}

func EventListenerPrefix(p string) EventListenerConfiguration {

	return func(e *EventListener) error {
		e.prefix = p
		return nil
	}
}

func EventListenerHTTPResponse(r *http.Response) EventListenerConfiguration {

	return func(e *EventListener) error {
		e.httpResponse = r
		return nil
	}
}

///////////////////////////////////////////////////////////////////////////
// Private methods
///////////////////////////////////////////////////////////////////////////

func (e *EventListener) url() (string, error) {

	url := &APIUrl{
		BaseUrl,
		APIVersion,
		"/devices",
	}

	// Add a device ID, if specified
	if e.Device != nil {
		url.Endpoint += "/" + e.Device.ID
	}

	// All events are at /events
	url.Endpoint += "/events"

	// Spark cloud also support prefixes
	if e.prefix != "" {
		url.Endpoint += "/" + e.prefix
	}

	// Get an auth token
	token, err := e.AuthTokenProvider.AuthToken()

	if err != nil {
		return "", err
	}

	// ... and add it
	url.Endpoint += "?access_token=" + token.String()

	return url.String(), nil
}

func (e *EventListener) connect() error {

	url, url_err := e.url()

	if url_err != nil {
		return url_err
	}

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("EventListener.connect(): response status code %d\n", resp.StatusCode)
	}

	// Store the response
	e.httpResponse = resp

	return nil
}

// Either return an error on the channel, or log it. Don't close
// the channels at this point, though; allow the host to do that.
func (e *EventListener) error(d string, a ...interface{}) error {

	err := fmt.Errorf(d, a...)

	if e.ErrorChan != nil {
		e.ErrorChan <- err
	} else {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	return err
}

///////////////////////////////////////////////////////////////////////////
// Public methods
///////////////////////////////////////////////////////////////////////////

func (e *EventListener) Listen() error {

	ev := Event{}

	reader := bufio.NewReader(e.httpResponse.Body)
	var buf bytes.Buffer

	e.running = true

	for {

		line, err := reader.ReadBytes('\n')

		if e.running == false {
			break
		}

		if err != nil {
			return e.error("error during resp.Body read:%s\n", err)
		}

		switch {

		// OK line. This is a keepalive; do nothing
		case bytes.HasPrefix(line, []byte(":ok")):

		// name of event
		case bytes.HasPrefix(line, []byte("event:")):
			ev.Name = string(line[7 : len(line)-1])

		// event data
		case bytes.HasPrefix(line, []byte("data:")):
			buf.Write(line[6:])

		// end of event
		case bytes.Equal(line, []byte("\n")):
			b := buf.Bytes()

			if bytes.HasPrefix(b, []byte("{")) {
				err := json.Unmarshal(b, &ev)

				if err == nil {
					buf.Reset()
					e.OutputChan <- ev
					ev = Event{}
				}
			}

		default:
			e.error("Error: len:%d\n%s", len(line), line)
		}
	}

	return nil
}

func (e *EventListener) Stop() {
	e.running = false
}

func (e *EventListener) Close() {
	close(e.ErrorChan)
	close(e.OutputChan)
	e.Stop()
}

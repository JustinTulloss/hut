package hut

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
)

// A Client is a convenient way of consuming other hut services
type Client interface {
	Serializer() Serializer
	Call(resource string, payload interface{}) (chan interface{}, chan error)
}

type BaseClient struct {
	serializer Serializer
}

func (c *BaseClient) Serializer() Serializer {
	return c.serializer
}

// Client for communicating with HTTP services
type HttpClient struct {
	BaseClient
	addresses []*url.URL
}

func (c *HttpClient) Call(resource string, payload interface{}) (chan *http.Response, chan error) {
	responseCh := make(chan *http.Response)
	errCh := make(chan error)
	go func() {
		buffer := bytes.NewBuffer(make([]byte, 1024))
		encoder := c.Serializer().NewEncoder(buffer)
		encoder.Encode(payload)
		body := bytes.NewReader(buffer.Bytes())
		for _, addr := range c.addresses {
			fullAddr := strings.Join([]string{addr.String(), resource}, "")
			resp, err := http.Post(fullAddr, c.Serializer().MimeType(), body)
			if err != nil {
				continue // Try the next address
			}
			responseCh <- resp
		}
	}()
	return responseCh, errCh
}

func NewHttpClient(serviceName string, s Serializer) (*HttpClient, error) {
	// First, discover all the services with this name
	address, err := url.Parse("localhost:4000")
	if err != nil {
		return nil, err
	}
	addresses := []*url.URL{address}
	return &HttpClient{BaseClient{serializer: s}, addresses}, nil
}

func NewJsonHttpClient(serviceName string) (*HttpClient, error) {
	return NewHttpClient(serviceName, &JsonSerializer{})
}

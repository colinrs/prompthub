package httpc

import "net/http"

type metadata struct {
	clientName  string
	header      http.Header
	handleError func(*http.Response, *http.Request) error
}

// Option ...
type Option func(*metadata)

// WithHeader ...
func WithHeader(header http.Header) Option {
	return func(c *metadata) {
		c.header = header
	}
}

// WithClientName ...
func WithClientName(clientName string) Option {
	return func(c *metadata) {
		c.clientName = clientName
	}
}

// NewMetadata ...
func newMetadata(options ...Option) *metadata {
	m := &metadata{
		header:      make(http.Header),
		clientName:  "default",
		handleError: handleError,
	}
	for _, option := range options {
		option(m)
	}
	return m
}

package httpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/rest/clientinterceptor"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

type Client interface {
	Get(ctx context.Context, path string, options ...Option) (*http.Response, error)
	Post(ctx context.Context, path string, body interface{}, options ...Option) (*http.Response, error)
	Patch(ctx context.Context, path string, body interface{}, options ...Option) (*http.Response, error)
	Delete(ctx context.Context, path string, options ...Option) (*http.Response, error)
}

// NewClient ...
func NewClient(domain string) Client {
	return &client{domain: domain}
}

type client struct {
	domain string
}

// Get Helper functions
func (s *client) Get(ctx context.Context, path string, options ...Option) (*http.Response, error) {
	return s.sendRequest(ctx, http.MethodGet, path, nil, options...)
}

func (s *client) Post(ctx context.Context, path string, body interface{}, options ...Option) (*http.Response, error) {
	return s.sendRequest(ctx, http.MethodPost, path, body, options...)
}

func (s *client) Patch(ctx context.Context, path string, body interface{}, options ...Option) (*http.Response, error) {
	return s.sendRequest(ctx, http.MethodPatch, path, body, options...)
}

func (s *client) Delete(ctx context.Context, path string, options ...Option) (*http.Response, error) {
	return s.sendRequest(ctx, http.MethodDelete, path, nil, options...)
}

func (s *client) sendRequest(ctx context.Context, method, path string, body interface{}, options ...Option) (*http.Response, error) {
	httpMetadata := newMetadata(options...)
	var jsonBody []byte
	var err error

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	url := s.domain + path
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	for k, v := range httpMetadata.header {
		request.Header.Set(k, v[0])
	}
	interceptor := clientinterceptor.MetricsInterceptor(httpMetadata.clientName, nil)
	request, handler := interceptor(request)
	resp, err := httpc.DoRequest(request)
	handler(resp, err)
	if err != nil {
		logx.WithContext(ctx).Errorf("Failed to send request:%+v err: %v", request, err)
		return nil, err
	}
	if err := httpMetadata.handleError(resp, request); err != nil {
		logx.WithContext(ctx).Errorf("Failed to send request:%+v,resp:%+v err: %v", request, resp, err)
		return nil, err
	}
	logx.WithContext(ctx).Infof("success call path:%s,request body:%+v, resp:%+v", path, body, resp)
	return resp, nil
}

func handleError(resp *http.Response, _ *http.Request) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return code.HTTPClientErr
}

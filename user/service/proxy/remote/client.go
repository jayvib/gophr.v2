package remote

import (
	"context"
	"gophr.v2/user"
	"io"
	"net/http"
	"net/url"
)

type Response struct {
	Data    *user.User `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
}

type ClientOpt func(client *Client) error

func newClientWithDefaults() (*Client, error) {
	return newClient(nil, SetBaseUrl(defaultUrl))
}

func newClient(baseClient *http.Client, opts ...ClientOpt) (*Client, error) {
	if baseClient == nil {
		baseClient = http.DefaultClient
	}

	client := &Client{
		client: baseClient,
	}

	for _, opt := range opts {
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

type Client struct {
	client    *http.Client
	baseURL   *url.URL
	UserAgent string
}

func (c *Client) NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method,u.String(), body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx != nil {
		req.WithContext(ctx)
	}
	return c.client.Do(req)
}

func SetBaseUrl(rawurl string) func(*Client) error {
	return func(client *Client) error {
		baseUrl, err := url.Parse(rawurl)
		if err != nil {
			return err
		}
		client.baseURL = baseUrl
		return nil
	}
}

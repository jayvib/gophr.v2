package remote

import (
	"gophr.v2/user"
	"net/http"
	"net/url"
)

type Response struct {
	Data    *user.User `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
}

type Client struct {
	client *http.Client
	BaseURL *url.URL
	UserAgent string
}

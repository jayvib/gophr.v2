package remote

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gophr.v2/user"
	"net/http"
	"net/url"
)

const (
	defaultUrl = "http://127.0.0.1:4401"
)

var _ user.Service = (*Service)(nil)

func New(client *http.Client) *Service {
	if client == nil {
		client = http.DefaultClient
	}

	baseUrl, err := url.Parse(defaultUrl)
	if err != nil {
		panic(err)
	}

	return &Service{
		client: client,
		baseUrl: baseUrl,
	}
}

type Service struct {
	client *http.Client
	baseUrl *url.URL
	user.Service
}

func (s *Service) GetByID(ctx context.Context, id interface{}) (*user.User, error) {
	u, err := s.baseUrl.Parse(fmt.Sprintf("/user/%v", id))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("not found")
	}
	defer resp.Body.Close()

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}


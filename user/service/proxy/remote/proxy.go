package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"gophr.v2/user"
	"io"
	"net/http"
)

const (
	defaultUrl = "http://127.0.0.1:4401"
)

var _ user.Service = (*Service)(nil)

func New(client *Client) *Service {
	return &Service{
		client: client,
	}
}

type Service struct {
	client *Client
	user.Service
}

func (s *Service) GetByID(ctx context.Context, id interface{}) (*user.User, error) {
	return s.doGet(ctx, fmt.Sprintf("/user/%v", id), nil)
}

func (s *Service) Register(ctx context.Context, user *user.User) error {
	// Marshal the user and create a byte reader
	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	body := bytes.NewReader(payload)

	// Create a request
	req, err := s.client.NewRequest(http.MethodPut, "/user", body)
	if err != nil {
		return err
	}

	response, err := s.client.Do(ctx, req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errRes Response
		err := json.NewDecoder(response.Body).Decode(&errRes)
		if err != nil {
			return err
		}
		return errors.New(errRes.Error)
	}

	var resp Response
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return err
	}

	if !resp.Success {
		return errors.New("failed registering user")
	}

	return copier.Copy(user, &resp.Data)
}

func (s *Service) doGet(ctx context.Context, path string, body io.Reader) (*user.User, error) {
	req, err := s.client.NewRequest(http.MethodGet, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, user.ErrNotFound
	case http.StatusInternalServerError:
		return nil, errors.New("unexpected error")
	}

	defer resp.Body.Close()

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (s *Service) GetByUserID(ctx context.Context, id string) (*user.User, error) {
	return s.doGet(ctx, fmt.Sprintf("/user/%v", id), nil)
}

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

func (s *Service) Update(ctx context.Context, usr *user.User) error {

	// Marshal to json
	payload, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	body := bytes.NewReader(payload)

	// Do Request
	req, err := s.client.NewRequest(http.MethodPost, "/user", body)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return err
	}
	defer noOpClose(resp.Body)

	err = s.checkErr(resp)
	if err != nil {
		return err
	}

	// Unmarshal the response
	var usrRes user.User
	err = json.NewDecoder(resp.Body).Decode(&usrRes)
	if err != nil {
		return err
	}

	return copier.Copy(usr, usrRes)
}

func (s *Service) Delete(ctx context.Context, id interface{}) error {

	req, err := s.client.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%s", id), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return err
	}
	defer noOpClose(resp.Body)

	err = s.checkErr(resp)
	if err != nil {
		return err
	}

	var r Response
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return err
	}

	if !r.Success {
		return errors.New("failed deleting an item")
	}
	return nil
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
	defer noOpClose(response.Body)

	err = s.checkErr(response)
	if err != nil {
		return err
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

func (s *Service) checkErr(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError:
		var errRes Response
		err := json.NewDecoder(resp.Body).Decode(&errRes)
		if err != nil {
			return err
		}
		return errors.New(errRes.Message)
	case http.StatusOK:
		return nil
	default:
		return errors.New(fmt.Sprintf("unhandled response status code: %d", resp.StatusCode))
	}
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

	defer noOpClose(resp.Body)

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

func noOpClose(c io.Closer) {
	_ = c.Close()
}
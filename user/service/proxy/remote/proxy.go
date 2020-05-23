package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/jinzhu/copier"
	"gophr.v2/user"
	"io"
	"net/http"
	"net/url"
	"reflect"
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
	return s.doGet(ctx, fmt.Sprintf("/users/%v", id), nil)
}

func (s *Service) Update(ctx context.Context, usr *user.User) error {

	// Marshal to json
	payload, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	body := bytes.NewReader(payload)

	// Do Request
	req, err := s.client.NewRequest(http.MethodPost, "/users", body)
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
	var requestResponse Response
	err = json.NewDecoder(resp.Body).Decode(&requestResponse)
	if err != nil {
		return err
	}

	usrRes, err := decodeUser(requestResponse)

	return copier.Copy(usr, usrRes)
}

func (s *Service) Delete(ctx context.Context, id interface{}) error {

	req, err := s.client.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", id), nil)
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
	req, err := s.client.NewRequest(http.MethodPut, "/users", body)
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

	usr, err := decodeUser(resp)
	if err != nil {
		return err
	}

	return copier.Copy(user, usr)
}

func decodeUser(resp Response) (*user.User, error) {
	payload, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}

	var usr user.User
	r := bytes.NewReader(payload)
	err = json.NewDecoder(r).Decode(&usr)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func (s *Service) GetByUserID(ctx context.Context, id string) (*user.User, error) {
	return s.doGet(ctx, fmt.Sprintf("/users/%v", id), nil)
}

func (s *Service) GetAll(ctx context.Context, cursor string, num int) (users []*user.User, next string, err error) {

	opt := &struct{
		Cursor string `url:"cursor,omitempty"`
		Num int `url:"num,omitempty"`
	} {
		 cursor,
		 num,
	}

	path, err := addOptions("/users", opt)
	if err != nil {
		return nil, "", err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, "", err
	}

	reqResp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, "", err
	}

	if err := s.checkErr(reqResp);err != nil {
	  return nil, "", err
	}

	var response Response
	err = json.NewDecoder(reqResp.Body).Decode(&response)
	if err != nil {
		return nil, "", err
	}

	users, err = decodeUsers(&response)
	if err != nil {
		return nil, "", err
	}

	// Get the cursor
	next = reqResp.Header.Get("X-Cursor")
	return
}

func decodeUsers(resp *Response) ([]*user.User, error) {
	payload, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}

	users := make([]*user.User, 0)
	r := bytes.NewReader(payload)
	err = json.NewDecoder(r).Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func addOptions(s string, option interface{}) (string, error) {
	v := reflect.ValueOf(option)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origUrl, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origUrl.Query()

	newValues, err := query.Values(option)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origUrl.RawQuery = origValues.Encode()
	return origUrl.String(), nil
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
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
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

	mapPayload, err := json.Marshal(result.Data)
	if err != nil {
		return nil, err
	}

	// unmarshal back for user.User
	var usr user.User

	err = json.Unmarshal(mapPayload, &usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func noOpClose(c io.Closer) {
	_ = c.Close()
}
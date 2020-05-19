package remote

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gophr.v2/user"
	"gophr.v2/user/userutil"
	"gophr.v2/util/valueutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetByUserID(t *testing.T) {
	t.Run("Found", func(t *testing.T){
		want := &user.User{
			ID: 1,
			UserID: userutil.GenerateID(),
			Email: "unit.test@testing.com",
			Password: "mysupersecretpassword",
		}
		response := &Response{
			Data: want,
			Success: true,
		}; _ = response

		// Create a mock handler
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			payload, err := json.Marshal(response)
			assert.NoError(t, err)
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(payload)
			assert.NoError(t, err)
		})

		c, teardown := setupClient(t, h)
		defer teardown()

		// Initialise the service
		svc := New(c)

		// Assert the result
		res, err := svc.GetByUserID(context.Background(), want.UserID)
		require.NoError(t, err)

		assert.Equal(t, want, res)
	})

	t.Run("Not Found", func(t *testing.T) {
		response := &Response{
			Success: false,
			Message: "Filed to get the user because it is not found",
		}

		// Create a mock handler
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			payload, err := json.Marshal(response)
			require.NoError(t, err)
			_, err = w.Write(payload)
			require.NoError(t, err)
		})

		c, teardown := setupClient(t, h)
		defer teardown()

		// Initialise the service
		svc := New(c)

		// Assert the result
		_, err := svc.GetByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Equal(t, user.ErrNotFound, err)
	})
}

func setupClient(t *testing.T, h http.HandlerFunc) (*Client, func()) {
	t.Helper()
	// Create an http client
	client, teardown := testingHTTPClient(h)
	c, err := newClient(client, SetBaseUrl(defaultUrl))
	require.NoError(t, err)
	return c, teardown
}

func TestRegister(t *testing.T) {

	t.Run("Success", func(t *testing.T){
		want := &user.User{
			Username: "unit.testing",
			Email: "unit.test@testing.com",
			Password: "mysupersecretpassword",
		}

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			var user user.User

			err := json.NewDecoder(r.Body).Decode(&user)
			require.NoError(t, err)

			user.UserID = userutil.GenerateID()

			password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			require.NoError(t, err)

			user.Password = string(password)


			response := &Response{
				Data: &user,
				Success: true,
			}
			payload, err := json.Marshal(response)
			assert.NoError(t, err)

			w.WriteHeader(http.StatusOK)
			w.Write(payload)
		})

		client, teardown := setupClient(t, h)
		defer teardown()

		svc := New(client)

		err := svc.Register(context.Background(), want)
		require.NoError(t, err)

		assert.NotEmpty(t, want.UserID)
		assert.NotEmpty(t, want.Password)
	})

	t.Run("Status Not OK", func(t *testing.T){
		want := &user.User{
			Username: "unit.testing",
			Email: "unit.test@testing.com",
			Password: "mysupersecretpassword",
		}

		wantMsg := "this is a unit test error"
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			response := &Response{
				Success: false,
				Message: wantMsg,
			}
			payload, err := json.Marshal(response)
			assert.NoError(t, err)

			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(payload)
		})

		client, teardown := setupClient(t, h)
		defer teardown()

		svc := New(client)
		err := svc.Register(context.Background(), want)
		assert.Error(t, err)
		assert.Equal(t, wantMsg, err.Error())
	})

}

func TestUpdate(t *testing.T) {
	want := &user.User{
		Username: "unit.testing",
		Email: "updated.unit.test@testing.com",
		Password: "mysupersecretpassword",
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	 var usr user.User
	 err := json.NewDecoder(r.Body).Decode(&usr)
	 require.NoError(t, err)
	 defer r.Body.Close()
	 usr.UpdatedAt = valueutil.TimePointer(time.Now().UTC())

		payload, err := json.Marshal(&usr)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	})

	client, teardown := setupClient(t, h)
	defer teardown()

	svc := New(client)

	err := svc.Update(context.Background(), want)
	assert.NoError(t, err)
	assert.NotEmpty(t, want.UpdatedAt)
}

func TestDelete(t *testing.T) {
}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)
	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}
	return cli, s.Close
}

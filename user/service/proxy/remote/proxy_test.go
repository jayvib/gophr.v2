package remote

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophr.v2/user"
	"gophr.v2/user/userutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
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

		// Create an http client
		client, teardown := testingHTTPClient(h)
		defer teardown()

		c, err := newClient(client, SetBaseUrl(defaultUrl))
		assert.NoError(t, err)

		// Initialise the service
		svc := New(c)

		// Assert the result
		res, err := svc.GetByUserID(context.Background(), want.UserID)
		require.NoError(t, err)

		assert.Equal(t, want, res)
	})

	t.Run("Not Found", func(t *testing.T){
		response := &Response{
			Success: false,
			Message: "Filed to get the user because it is not found",
		}

		// Create a mock handler
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			w.WriteHeader(http.StatusNotFound)
			payload, err := json.Marshal(response)
			require.NoError(t, err)
			_, err = w.Write(payload)
			require.NoError(t, err)
		})

		// Create an http client
		client, teardown := testingHTTPClient(h)
		defer teardown()

		c, err := newClient(client, SetBaseUrl(defaultUrl))
		assert.NoError(t, err)

		// Initialise the service
		svc := New(c)

		// Assert the result
		_, err = svc.GetByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Equal(t, user.ErrNotFound, err)
	})
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

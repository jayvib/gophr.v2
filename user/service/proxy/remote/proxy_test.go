package remote

import (
	"context"
	"encoding/json"
	"github.com/jayvib/golog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophr.v2/user"
	"gophr.v2/user/userutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestService_GetByID(t *testing.T) {
	golog.SetLevel(golog.DebugLevel)
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

	// Initialise the service
	svc := New(client)

	// Assert the result
	res, err := svc.GetByID(context.Background(), 1)
	require.NoError(t, err)

	assert.Equal(t, want, res)
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

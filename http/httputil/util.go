package httputil

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
)

func PerformRequest(h http.Handler, method string, path string, body io.Reader, opts ...func(r *http.Request)) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	for _, opt := range opts {
		opt(req)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func DummyClient(h http.Handler) (cli *http.Client, teardown func()) {
	s := httptest.NewServer(h)
	cli = &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}
	return cli, s.Close
}

package httputil

import (
  "io"
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

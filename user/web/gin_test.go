//+build unit

package web

import (
  "encoding/json"
  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "gophr.v2/errors"
  "gophr.v2/user"
  "gophr.v2/user/mocks"
  "gophr.v2/user/service"
  "io"
  "net/http"
  "net/http/httptest"
  "os"
  "testing"
)

func TestMain(m *testing.M) {
  gin.SetMode(gin.TestMode)
  os.Exit(m.Run())
}

func TestGetByID(t *testing.T) {
  t.Run("StatusOK", func(t *testing.T){
    want := &user.User{
      ID: 1,
      Username: "luffy.monkey",
      Email: "luffy.monkey@gmail.com",
    }

    e := gin.Default()
    repo := new(mocks.Repository)
    repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(want, nil)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := performRequest(e, http.MethodGet, "/users/1",nil)

    require.Equal(t, http.StatusOK, response.Code)
    assertGetByID(t, want, response.Body)
  })

  t.Run("NotFound", func(t *testing.T){
    e := gin.Default()
    repo := new(mocks.Repository)
    repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.ErrorNotFound)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := performRequest(e, http.MethodGet, "/users/1",nil)

    assert.Equal(t, http.StatusNotFound, response.Code)
  })
}

func assertGetByID(t *testing.T, want *user.User, body io.Reader) {
  var got user.User
  err := json.NewDecoder(body).Decode(&got)
  assert.NoError(t, err)
  assert.Equal(t, want, &got)
}

func performRequest(h http.Handler, method string, path string, body io.Reader) *httptest.ResponseRecorder {
  req := httptest.NewRequest(method, path, body)
  w := httptest.NewRecorder()
  h.ServeHTTP(w, req)
  return w
}

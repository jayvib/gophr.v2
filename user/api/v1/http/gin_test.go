//+build unit

package http

import (
  "encoding/json"
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/jayvib/golog"
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

var debug = flag.Bool("debug", false, "Debugging")

func TestMain(m *testing.M) {
  flag.Parse()
  if *debug {
    golog.Info("Debugging Mode!")
    golog.SetLevel(golog.DebugLevel)
  }
  gin.SetMode(gin.TestMode)
  os.Exit(m.Run())
}

func TestGetByID(t *testing.T) {
  t.Run("StatusOK", func(t *testing.T){
    usr := &user.User{
      ID: 1,
      Username: "luffy.monkey",
      Email: "luffy.monkey@gmail.com",
    }

    want := Response{
      Success: true,
      Data: usr,
    }

    e := gin.Default()
    repo := new(mocks.Repository)
    repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(usr, nil)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := performRequest(e, http.MethodGet, "/users/id/1",nil)

    require.Equal(t, http.StatusOK, response.Code)
    assertGetByID(t, want, response.Body)
    repo.AssertExpectations(t)
  })

  t.Run("NotFound", func(t *testing.T){
    e := gin.Default()
    want := Response{
      Success: false,
      Error: "item not found",
    }
    repo := new(mocks.Repository)
    repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.ErrorNotFound)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := performRequest(e, http.MethodGet, "/users/id/1",nil)

    assert.Equal(t, http.StatusNotFound, response.Code)
    assertGetByID(t, want, response.Body)
  })
}

func TestGetByEmail(t *testing.T) {
  t.Run("StatusOK", func(t *testing.T){
    usr := &user.User{
      ID: 1,
      Username: "luffy.monkey",
      Email: "luffy.monkey@gmail.com",
    }

    want := Response{
      Success: true,
      Data: usr,
    }

    e := gin.Default()
    repo := new(mocks.Repository)
    repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(usr, nil)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := performRequest(e, http.MethodGet, "/users/email/luffy.monkey@gmail.com",nil)

    require.Equal(t, http.StatusOK, response.Code)
    assertGetByID(t, want, response.Body)
    repo.AssertExpectations(t)
  })

  t.Run("NotFound", func(t *testing.T){
    t.SkipNow()
    e := gin.Default()
    want := Response{
      Success: false,
      Error: "item not found",
    }
    repo := new(mocks.Repository)
    repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := performRequest(e, http.MethodGet, "/users/email/luffy.monkey@gmail.com",nil)

    assert.Equal(t, http.StatusNotFound, response.Code)
    assertGetByID(t, want, response.Body)
  })
}

func assertGetByID(t *testing.T, want Response, body io.Reader) {
  t.Helper()
  var got Response
  err := json.NewDecoder(body).Decode(&got)
  assert.NoError(t, err)

  err, gotUser := extractUserFromData(err, got)
  require.NoError(t, err)

  if gotUser != nil {
    got.Data = gotUser
  }

  assert.Equal(t, want, got)
}

func extractUserFromData(err error, got Response) (error, *user.User) {
  if got.Data == nil {
    return nil, nil
  }
  usrPayload, err := json.Marshal(got.Data)
  var gotUser user.User
  err = json.Unmarshal(usrPayload, &gotUser)
  return err, &gotUser
}



func performRequest(h http.Handler, method string, path string, body io.Reader) *httptest.ResponseRecorder {
  req := httptest.NewRequest(method, path, body)
  w := httptest.NewRecorder()
  h.ServeHTTP(w, req)
  return w
}

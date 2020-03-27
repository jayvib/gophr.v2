//+build unit

package http

import (
  "bytes"
  "encoding/json"
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/jayvib/golog"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
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
    assertResponse(t, want, response.Body)
    repo.AssertExpectations(t)
  })

  t.Run("NotFound", func(t *testing.T){
    e := gin.Default()
    want := Response{
      Success: false,
      Error: "user: item not found",
    }
    repo := new(mocks.Repository)
    repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := performRequest(e, http.MethodGet, "/users/id/1",nil)

    assert.Equal(t, http.StatusNotFound, response.Code)
    assertResponse(t, want, response.Body)
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
    assertResponse(t, want, response.Body)
    repo.AssertExpectations(t)
  })

  t.Run("NotFound", func(t *testing.T){
    e := gin.Default()
    want := Response{
      Success: false,
      Error: "user: item not found",
    }
    repo := new(mocks.Repository)
    repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := performRequest(e, http.MethodGet, "/users/email/luffy.monkey@gmail.com",nil)

    assert.Equal(t, http.StatusNotFound, response.Code)
    assertResponse(t, want, response.Body)
  })
}

func TestGetByUsername(t *testing.T) {
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
    repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(usr, nil)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := performRequest(e, http.MethodGet, "/users/username/luffy.monkey",nil)

    require.Equal(t, http.StatusOK, response.Code)
    assertResponse(t, want, response.Body)
    repo.AssertExpectations(t)
  })

  t.Run("NotFound", func(t *testing.T){
    e := gin.Default()
    want := Response{
      Success: false,
      Error: "user: item not found",
    }
    repo := new(mocks.Repository)
    repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := performRequest(e, http.MethodGet, "/users/email/luffy.monkey@gmail.com",nil)

    assert.Equal(t, http.StatusNotFound, response.Code)
    assertResponse(t, want, response.Body)
  })
}

func TestRegister(t *testing.T) {
  t.Run("StatusCreated", func(t *testing.T){
    usr := &user.User{
      Username: "luffy.monkey",
      Email: "luffy.monkey@gmail.com",
      Password: "iampirateking",
    }
    e := gin.Default()
    repo := new(mocks.Repository)
    repo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
    repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound).Once()


    svc := service.New(repo)

    RegisterHandlers(e, svc)

    payload, err := json.Marshal(usr)
    require.NoError(t, err)

    body := bytes.NewReader(payload)
    response := performRequest(e, http.MethodPost, "/users", body)

    assert.Equal(t, http.StatusCreated, response.Code)
    var got Response
    err = json.NewDecoder(response.Body).Decode(&got)
    require.NoError(t, err)
    assert.True(t, got.Success)

    gotUser, err := extractUserFromData(got)
    require.NoError(t, err)
    assert.NotEmpty(t, gotUser.CreatedAt)
    assert.NotEmpty(t, gotUser.Password)

    repo.AssertExpectations(t)
  })

  t.Run("Validate User No Username", func(t *testing.T){
    usr := &user.User{
      Email: "luffy.monkey@gmail.com",
      Password: "iampirateking",
    }

    e := gin.Default()
    svc := service.New(nil)
    RegisterHandlers(e, svc)

    payload, err := json.Marshal(usr)
    require.NoError(t, err)

    body := bytes.NewReader(payload)
    response := performRequest(e, http.MethodPost, "/users", body)
    assert.Equal(t, http.StatusBadRequest, response.Code)
  })
}

func assertResponse(t *testing.T, want Response, body io.Reader) {
  t.Helper()
  var got Response
  err := json.NewDecoder(body).Decode(&got)
  assert.NoError(t, err)

  gotUser, err := extractUserFromData(got)
  require.NoError(t, err)

  if gotUser != nil {
    got.Data = gotUser
  }

  assert.Equal(t, want, got)
}

func extractUserFromData(got Response) (*user.User, error) {
  if got.Data == nil {
    return nil, nil
  }
  usrPayload, err := json.Marshal(got.Data)
  var gotUser user.User
  err = json.Unmarshal(usrPayload, &gotUser)
  return &gotUser, err
}



func performRequest(h http.Handler, method string, path string, body io.Reader) *httptest.ResponseRecorder {
  req := httptest.NewRequest(method, path, body)
  w := httptest.NewRecorder()
  h.ServeHTTP(w, req)
  return w
}

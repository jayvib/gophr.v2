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
  "gophr.v2/http/httputil"
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

    response := httputil.PerformRequest(e, http.MethodGet, "/users/id/1",nil)

    require.Equal(t, http.StatusOK, response.Code)
    assertResponse(t, want, response.Body)
    repo.AssertExpectations(t)
  })

  t.Run("NotFound", func(t *testing.T){
    e := gin.Default()
    want := Response{
      Success: false,
      Message: "Failed getting the user because it didn't exist",
    }
    repo := new(mocks.Repository)
    repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound)

    svc := service.New(repo)
    RegisterHandlers(e, svc)

    response := httputil.PerformRequest(e, http.MethodGet, "/users/id/1",nil)

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

    body := userToBody(t, usr)
    response := httputil.PerformRequest(e, http.MethodPut, "/users", body)

    assert.Equal(t, http.StatusCreated, response.Code)
    got := extractResponse(t, response)
    assert.True(t, got.Success)
    gotUser, err := extractUserFromData(got)
    require.NoError(t, err)
    assert.NotEmpty(t, gotUser.CreatedAt)
    assert.NotEmpty(t, gotUser.Password)

    repo.AssertExpectations(t)
  })

  t.Run("Validate User No Username", func(t *testing.T) {
    usr := &user.User{
      Email:    "luffy.monkey@gmail.com",
      Password: "iampirateking",
    }

    e := gin.Default()
    svc := service.New(nil)
    RegisterHandlers(e, svc)

    body := userToBody(t, usr)
    response := httputil.PerformRequest(e, http.MethodPut, "/users", body)
    got := extractResponse(t, response)
    assert.Equal(t, http.StatusBadRequest, response.Code)
    assert.Equal(t, "Missing value for:\nUsername\n", got.Message)
  })
}

func TestDelete(t *testing.T) {
  e := gin.Default()
  svc := new(mocks.Service)
  svc.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
  RegisterHandlers(e, svc)
  response := httputil.PerformRequest(e, http.MethodDelete, "/users/id/:id", nil)
  assert.Equal(t, http.StatusOK, response.Code)
  got := extractResponse(t, response)
  assert.True(t, got.Success)
  svc.AssertExpectations(t)
}


func TestUpdate(t *testing.T) {
  t.Run("When Updating an Existed User", func(t *testing.T){
    repo := new(mocks.Repository)
    repo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(nil, nil).Once()
    repo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
    svc := service.New(repo)

    e := gin.Default()
    RegisterHandlers(e, svc)

    input := &user.User{
      ID: 1,
      Username: "luffy.monkey",
      Email: "luffy.monkey@gmail.com",
      Password: "Secret",
    }

    body := userToBody(t, input)
    response := httputil.PerformRequest(e, http.MethodPost, "/users", body)
    assert.Equal(t, http.StatusOK, response.Code)
    repo.AssertExpectations(t)

    got := extractResponse(t, response)
    assert.True(t, got.Success)
  })

  t.Run("When Updating to an Non-Exiting User", func(t *testing.T){
    repo := new(mocks.Repository)
    repo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(nil, user.ErrNotFound).Once()

    svc := service.New(repo)

    e := gin.Default()
    RegisterHandlers(e, svc)

    input := &user.User{
      ID: 1,
      Username: "luffy.monkey",
      Email: "luffy.monkey@gmail.com",
      Password: "Secret",
    }

    body := userToBody(t, input)
    response := httputil.PerformRequest(e, http.MethodPost, "/users", body)
    assert.Equal(t, http.StatusBadRequest, response.Code)
    repo.AssertExpectations(t)

    got := extractResponse(t, response)
    assert.False(t, got.Success)
    assert.Equal(t, "Failed because user is not exists", got.Message)
  })
}

func userToBody(t *testing.T, usr *user.User) *bytes.Reader {
  payload, err := json.Marshal(usr)
  require.NoError(t, err)
  body := bytes.NewReader(payload)
  return body
}

func extractResponse(t *testing.T, response *httptest.ResponseRecorder) Response {
  t.Helper()
  var got Response
  err := json.NewDecoder(response.Body).Decode(&got)
  require.NoError(t, err)
  return got
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




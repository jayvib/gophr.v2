package web

import (
  "encoding/json"
  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "gophr.v2/user"
  "gophr.v2/user/mocks"
  "gophr.v2/user/service"
  "net/http"
  "os"
  "testing"
)

func TestMain(m *testing.M) {
  gin.SetMode(gin.TestMode)
  os.Exit(m.Run())
}

func TestGetByID(t *testing.T) {
  want := &user.User{
    ID: 1,
    Username: "luffy.monkey",
    Email: "luffy.monkey@gmail.com",
  }

  e := gin.Default()
  repo := new(mocks.Repository)
  repo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(want, nil)

  svc := service.New(repo)
  h := &GinHandler{svc:svc}
  e.GET("/users/:id", h.GetByID)
  response := performRequest(e, http.MethodGet, "/users/1",nil)

  require.Equal(t, http.StatusOK, response.Code)

  var got user.User
  err := json.NewDecoder(response.Body).Decode(&got)
  assert.NoError(t, err)
  assert.Equal(t, want, got)
}


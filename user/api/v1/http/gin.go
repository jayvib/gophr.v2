package http

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/go-playground/validator/v10"
  "github.com/jayvib/golog"
	"gophr.v2/user"
	"net/http"
  "strconv"
  "strings"
)

var validater = validator.New()

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
}

func RegisterHandlers(r gin.IRouter, svc user.Service) {
	handler := New(svc)
	r.GET("/users/id/:id", handler.GetByID)
	r.GET("/users", handler.GetAll)
  r.PUT("/users", handler.Register)
	r.POST("/users", handler.Update)
  r.DELETE("/users/id/:id", handler.Delete)
}

func New(svc user.Service) *GinHandler {
	return &GinHandler{
		svc: svc,
	}
}

type GinHandler struct {
	svc user.Service
}

func (g *GinHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	g.get(c, id, g.svc.GetByID)
}

func (g *GinHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")
	golog.Debug("Email:", email)
	g.get(c, email, g.svc.GetByEmail)
}

func (g *GinHandler) GetByUsername(c *gin.Context) {
  username := c.Param("username")
  g.get(c, username, g.svc.GetByUsername)
}

func (g *GinHandler) Delete(c *gin.Context) {
  id := c.Param("id")
  err := g.svc.Delete(c.Request.Context(), id)
  if err != nil {
    g.renderError(c, err)
    return
  }
  g.renderData(c, http.StatusOK, nil)
}
func (g *GinHandler) Update(c *gin.Context) {
  usr, err := g.decodeUserFromBody(c)
  if err != nil {
    golog.Debug("error:", err)
    g.renderError(c, err)
    return
  }

  err = g.svc.Update(c.Request.Context(), usr)
  if err != nil {
    golog.Debug("error:", err)
    g.renderError(c, err)
    return
  }

  g.renderData(c, http.StatusOK, usr)
}

func (g *GinHandler) Register(c *gin.Context) {
  usr, err := g.decodeUserFromBody(c)
  if err != nil {
    g.renderError(c, err)
    return
  }

  err = validater.Struct(usr)
  if err != nil {
    golog.Debugf("%T\n", err)
    g.renderError(c, err)
    return
  }

  err = g.svc.Register(c.Request.Context(), usr)
  if err != nil {
    golog.Debug(err.Error())
    g.renderError(c, err)
    return
  }
  g.renderData(c, http.StatusCreated, usr)
}

func (g *GinHandler) GetAll(c *gin.Context) {
  numString := c.Query("num")
  num, _ := strconv.Atoi(numString)
  cursor := c.Query("cursor")

  usrs, nextCursor, err := g.svc.GetAll(c.Request.Context(), cursor, num)
  if err != nil {
    g.renderError(c, err)
    return
  }

  c.Header(`X-Cursor`, nextCursor)
  g.renderData(c, http.StatusOK, usrs)
}

func (g *GinHandler) Login(c *gin.Context) {}

func (g *GinHandler) decodeUserFromBody(c *gin.Context) (*user.User, error) {
  var usr user.User
  err := json.NewDecoder(c.Request.Body).Decode(&usr)
  if err != nil {
    return nil, err
  }
  return &usr, err
}

func getStatusFromError(err error) int {
  golog.Debug(err)
  var status int
  switch errors.Unwrap(err) {
  case user.ErrEmptyUsername, user.ErrEmptyEmail, user.ErrEmptyPassword, user.ErrUserExists:
    status = http.StatusBadRequest
  case user.ErrUserNotExists:
    status = http.StatusBadRequest
  case user.ErrNotFound:
    status = http.StatusNotFound
  default:
    switch err.(type) {
    case validator.ValidationErrors, *json.SyntaxError:
      status = http.StatusBadRequest
    default:
      status = http.StatusInternalServerError
    }
  }
  return status
}

func generateMessageFromError(err error) string {
  switch e := err.(type) {
  case validator.ValidationErrors:
    var b strings.Builder
    _, _ = fmt.Fprint(&b, "Missing value for:\n")
    for _, verr := range e {
      fieldName := verr.StructField()
      _, _ = fmt.Fprintln(&b, fieldName)
    }
    return b.String()
  case *user.Error:
    return e.Message()
  default:
    switch err {
    case user.ErrUserExists:
      return "Update user failed because it did not exists"
    default:
      return ""
    }
  }
}

func (g *GinHandler) renderError(c *gin.Context, err error) {
  if uerr, ok := err.(*user.Error); ok {
    golog.Error(uerr)
  }

  c.JSON(getStatusFromError(err), &Response{
    Success: false,
    Message: generateMessageFromError(err),
  })
}

func (g *GinHandler) renderData(c *gin.Context, status int, data interface{}) {
  c.JSON(status, &Response{
    Success: true,
    Data: data,
  })
}

func (g *GinHandler) get(c *gin.Context, id interface{}, getterFunc interface{}) {
  var usr *user.User
  var err error

  switch fn := getterFunc.(type) {
  case func(ctx context.Context, id interface{})(*user.User, error):
    usr, err = fn(c.Request.Context(), id)
  case func(ctx context.Context, input string)(*user.User, error):
    usr, err = fn(c.Request.Context(), id.(string))
  }

  if err != nil {
    g.renderError(c, err)
    return
  }

  g.renderData(c, http.StatusOK, usr)
}

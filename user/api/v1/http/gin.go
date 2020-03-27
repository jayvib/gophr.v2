package http

import (
  "context"
  "github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	"gophr.v2/user"
	"net/http"
)

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success,omitempty"`
	Message string      `json:"message,omitempty"`
}

func RegisterHandlers(r gin.IRouter, svc user.Service) {
	handler := New(svc)
	r.GET("/users/id/:id", handler.GetByID)
	r.GET("/users/email/:email", handler.GetByEmail)
	r.GET("/users/username/:username", handler.GetByUsername)
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
    c.JSON(http.StatusNotFound,
      Response{
        Error:   err.Error(),
        Success: false,
      })
    return
  }

  c.JSON(http.StatusOK, &Response{
    Data:    usr,
    Success: true,
  })
}

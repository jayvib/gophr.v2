package web

import (
  "github.com/gin-gonic/gin"
  "gophr.v2/user"
  "net/http"
)


func RegisterHandlers(r gin.IRouter, svc user.Service) {
  handler := New(svc)
  r.GET("/users/:id", handler.GetByID)
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

  usr, err := g.svc.GetByID(c.Request.Context(), id)
  if err != nil {
    c.Writer.WriteHeader(http.StatusNotFound)
    return
  }
  c.JSON(http.StatusOK, usr)
}

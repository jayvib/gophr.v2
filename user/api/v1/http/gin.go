package http

import (
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

func (g *GinHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")
	golog.Debug("Email:", email)

	usr, err := g.svc.GetByEmail(c.Request.Context(), email)
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

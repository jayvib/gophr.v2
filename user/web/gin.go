package web

import (
  "github.com/gin-gonic/gin"
  "gophr.v2/user"
)

type GinHandler struct {
  svc user.Service
}

func (g *GinHandler) GetByID(c *gin.Context) {}

package v1

import (
  "github.com/gin-gonic/gin"
  "gophr.v2/image"
)

type handlers struct {
  svc image.Service
}

func (h *handlers) Save(c *gin.Context) {}
func (h *handlers) Find(c *gin.Context) {}
func (h *handlers) FindAll(c *gin.Context) {}
func (h *handlers) FindAllByUser(c gin.Context) {}
func (h *handlers) CreateImageFromURL(c gin.Context) {}
func (h *handlers) CreateImageFromFile(c gin.Context) {}

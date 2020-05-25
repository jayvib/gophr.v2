package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	"gophr.v2/image"
	"gophr.v2/user"
	"net/http"
)

func RegisterRoutes(r gin.IRouter, imageSvc image.Service, userSvc user.Service) {

	h := handlers{
		imageSvc: imageSvc,
		userSvc:  userSvc,
	}

	r.POST("/image/file", h.CreateImageFromFile)
	r.POST("/image/url", h.CreateImageFromURL)
}

type handlers struct {
	imageSvc image.Service
	userSvc  user.Service
}

func (h *handlers) Find(c *gin.Context)               {}
func (h *handlers) FindAll(c *gin.Context)            {}
func (h *handlers) FindAllByUser(c *gin.Context)      {}

func (h *handlers) CreateImageFromURL(c *gin.Context) {

	desc := c.PostForm("description")
	userName := c.PostForm("username")
	url := c.PostForm("url")

	usr, err := h.userSvc.GetByUsername(c.Request.Context(), userName)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		golog.Error("unable to create image:", err)
		return
	}

	img, err := h.imageSvc.CreateImageFromURL(c.Request.Context(), url, usr.UserID, desc)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		golog.Error("unable to create image from url:", err)
		return
	}

	c.JSON(http.StatusCreated, img)
}

func (h *handlers) CreateImageFromFile(c *gin.Context) {

	desc := c.PostForm("description")
	userName := c.PostForm("username")

	// Get the username
	usr, err := h.userSvc.GetByUsername(c.Request.Context(), userName)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		golog.Error("unable to create image:", err)
		return
	}

	formFile, err := c.FormFile("file")
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		golog.Error("file not exist in the multipart form:", err)
		return
	}

	f, err := formFile.Open()
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		golog.Error("unable to open form file:", err)
		return
	}
	defer f.Close()

	img, err := h.imageSvc.CreateImageFromFile(c.Request.Context(), f, formFile.Filename, desc, usr.UserID)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		golog.Error("failed to create image from file:", err)
	}

	c.JSON(http.StatusCreated, img)
}

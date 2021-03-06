package view

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	"golang.org/x/crypto/bcrypt"
	"gophr.v2/image"
	"gophr.v2/session"
	"gophr.v2/session/sessionutil"
	"gophr.v2/user"
	"html/template"
	"net/http"
)

var funcs = template.FuncMap{
	"yield": func() (string, error) {
		return "", nil
	},
	"navbar": func() (string, error) {
		return "", nil
	},
}

func RegisterRoutes(unsecuredRouter, securedRouter gin.IRoutes, userService user.Service, sessionService session.Service, imageService image.Service, templatesGlob, layoutPath, assetsPath, imagesPath string) {
	h := NewHandler(userService, sessionService, imageService, templatesGlob, layoutPath)

	// Asset handler
	unsecuredRouter.StaticFS("/assets", http.Dir(assetsPath))
	unsecuredRouter.StaticFS("/im/", http.Dir(imagesPath))
	unsecuredRouter.GET("/", h.HomePage)
	unsecuredRouter.GET("/signup", h.SignupPage)
	unsecuredRouter.GET("/login", h.LoginPage)
	unsecuredRouter.POST("/signup", h.HandleSignUp)
	unsecuredRouter.POST("/login", h.HandleLogin)

	securedRouter.GET("/account", h.EditUserPage)
	securedRouter.GET("/signout", h.SignOutPage)
	securedRouter.GET("/images/new", h.UploadImagePage)
	securedRouter.GET("/images/id/:imageID", h.ShowImage)
	securedRouter.POST("/account", h.HandleEditUser)
	securedRouter.POST("/images/new", h.HandleImageUpload)
}

func NewHandler(userService user.Service, sessionService session.Service, imageService image.Service, templatesGlob, layoutPath string) *ViewHandler {
	return &ViewHandler{
		usrService:     userService,
		sessionService: sessionService,
		imageService:   imageService,
		templs:         template.Must(template.ParseGlob(templatesGlob)),
		layout: template.Must(template.New("layout.html").
			Funcs(funcs).ParseFiles(layoutPath)),
	}
}

type ViewHandler struct {
	templs         *template.Template
	layout         *template.Template
	usrService     user.Service
	sessionService session.Service
	imageService   image.Service
}

// #################CONTROLLERS################
func (v *ViewHandler) HandleSignUp(c *gin.Context) {
	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")
	usr := &user.User{
		Email:    email,
		Username: username,
		Password: password,
	}

	err := v.usrService.Register(c.Request.Context(), usr)
	if err != nil {
		v.renderTemplate(c, "users/signup", map[string]interface{}{
			"User":  usr,
			"Error": getMessage(err),
		})
		return
	}

	sess := sessionutil.WriteSessionTo(c.Writer)
	sess.UserID = usr.UserID

	err = v.sessionService.Save(c.Request.Context(), sess)
	if err != nil {
		v.renderErrorTemplate(c, err)
		return
	}

	c.Redirect(http.StatusFound, "/?flash=User+created")
}

func (v *ViewHandler) HandleImageUpload(c *gin.Context) {
	switch {
	case c.PostForm("url") != "":
		v.createImageFromURL(c)
	default:
		v.createImageFromFile(c)
	}
}

func (v *ViewHandler) createImageFromURL(c *gin.Context) {
	// Get the user from the session
	usr := v.getUserFromCookie(c)
	url := c.PostForm("url")
	desc := c.PostForm("description")

	// Create an image object
	img, err := v.imageService.CreateImageFromURL(c.Request.Context(), url, usr.UserID, desc)
	if err != nil {
		v.renderTemplate(c, "images/new", map[string]interface{}{
			"Error":    err,
			"ImageURL": url,
			"Image":    img,
		})
		return
	}
	c.Redirect(http.StatusFound, "/?flash=Image+Uploaded+Successfully")
}

func (v *ViewHandler) createImageFromFile(c *gin.Context) {
	usr := v.getUserFromCookie(c)
	desc := c.PostForm("description")
	formFile, err := c.FormFile("file")
	if err != nil {
		v.renderTemplate(c, "images/new", map[string]interface{}{
			"Error": err,
		})
		return
	}
	f, err := formFile.Open()
	if err != nil {
		v.renderTemplate(c, "images/new", map[string]interface{}{
			"Error": err,
		})
		return
	}
	defer func() {
		_ = f.Close()
	}()

	img, err := v.imageService.CreateImageFromFile(c.Request.Context(), f, formFile.Filename, desc, usr.UserID)
	if err != nil {
		golog.Error(err)
		v.renderTemplate(c, "images/new", map[string]interface{}{
			"Error": err,
			"Image": img,
		})
		return
	}
	c.Redirect(http.StatusFound, "/?flash=Image+Uploaded+Successfully")
}

func (v *ViewHandler) HandleLogin(c *gin.Context) {
	// Get the credentials
	username := c.PostForm("username")
	password := c.PostForm("password")
	next := c.PostForm("next")
	usr := &user.User{
		Username: username,
		Password: password,
	}
	golog.Debug("password:", password)
	// Get the user detail through username
	err := v.usrService.Login(c.Request.Context(), usr)
	if err != nil {
		v.renderTemplate(c, "sessions/login", map[string]interface{}{
			"Error": getMessage(err),
			"User":  usr,
			"Next":  next,
		})
		return
	}

	// Create a session
	sess := sessionutil.WriteSessionTo(c.Writer)
	sess.UserID = usr.UserID

	// Save the session

	err = v.sessionService.Save(c.Request.Context(), sess)
	if err != nil {
		v.renderErrorTemplate(c, err)
		return
	}

	// When next is empty string then set it to '/' as default
	if next == "" {
		next = "/"
	}

	c.Redirect(http.StatusFound, next+"?flash=Signed+in")
}

func (v *ViewHandler) HandleEditUser(c *gin.Context) {
	// Get the current user
	usr := v.getUserFromCookie(c)

	// Get the updated information
	email := c.PostForm("email")
	newPassword := c.PostForm("newPassword")
	currentPassword := c.PostForm("currentPassword")

	// Check first if the current password is correct

	tmpUser := &user.User{
		Email:    email,
		Username: usr.Username,
		Password: currentPassword,
	}

	if newPassword != "" {
		err := v.usrService.Login(c.Request.Context(), tmpUser)
		if err != nil {
			v.renderTemplate(c, "users/edit", map[string]interface{}{
				"Error": getMessage(err),
				"User":  tmpUser,
			})
			return
		}

		// ===============NOT SURE IF IT BELONGS HERE==============
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			err := user.NewError(err)
			v.renderTemplate(c, "users/edit", map[string]interface{}{
				"Error": getMessage(err),
				"User":  tmpUser,
			})
		}
		usr.Password = string(hash)
		// ========================================================
	}

	usr.Email = email

	golog.Tracef("Updating: %#v\n", usr)
	// Save to repository
	err := v.usrService.Update(c.Request.Context(), usr)
	if err != nil {
		v.renderTemplate(c, "users/edit", map[string]interface{}{
			"Error": getMessage(err),
			"User":  tmpUser,
		})
	}

	// redirect to "/account"
	c.Redirect(http.StatusFound, "/account?flash=User+updated")
}

func (v *ViewHandler) renderErrorTemplate(c *gin.Context, err error) {
	v.renderTemplate(c, "other/error", map[string]interface{}{
		"Error": getMessage(err),
	})
}

func getMessage(err error) string {
	type messenger interface {
		Message() string
	}
	var message string
	if msger, ok := err.(messenger); ok {
		message = msger.Message()
	} else {
		message = err.Error()
	}
	return message
}

func (v *ViewHandler) HomePage(c *gin.Context) {
	images, err := v.imageService.FindAll(c.Request.Context(), 0)
	if err != nil {
		v.renderTemplate(c, "index/home", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	v.renderTemplate(c, "index/home", map[string]interface{}{
		"Images": images,
	})
}

// ###################VIEW####################

func (v *ViewHandler) SignupPage(c *gin.Context) {
	v.renderTemplate(c, "users/signup", nil)
}

func (v *ViewHandler) LoginPage(c *gin.Context) {
	next := c.Query("next")
	v.renderTemplate(c, "sessions/login", map[string]interface{}{
		"Next": next,
	})
}

func (v *ViewHandler) EditUserPage(c *gin.Context) {
	usr := v.getUserFromCookie(c)
	v.renderTemplate(c, "users/edit", map[string]interface{}{
		"User":        usr,
		"CurrentUser": usr,
	})
}

func (v *ViewHandler) SignOutPage(c *gin.Context) {
	// Get session
	sess := v.getSessionFromRequest(c)

	// Delete session if not empty
	if sess != nil {
		_ = v.sessionService.Delete(c.Request.Context(), sess.ID)
	}

	// Render the signout template
	v.renderTemplate(c, "sessions/signout", nil)
}

func (v *ViewHandler) UploadImagePage(c *gin.Context) {
	v.renderTemplate(c, "images/new", nil)
}

func (v *ViewHandler) UserEditPage(c *gin.Context) {}

func (v *ViewHandler) DisplayUserDetails(c *gin.Context) {}

func (v *ViewHandler) ShowImage(c *gin.Context) {
	// Get image by image ID
	imageId := c.Param("imageID")
	img, err := v.imageService.Find(c.Request.Context(), imageId)
	if err != nil {
		v.renderErrorTemplate(c, err)
		return
	}

	// Find user by user ID
	usr, err := v.usrService.GetByUserID(c.Request.Context(), img.UserID)
	if err != nil {
		v.renderErrorTemplate(c, err)
		return
	}

	// Render template
	v.renderTemplate(c, "images/show", map[string]interface{}{
		"Image": img,
		"User":  usr,
	})

}

func (v *ViewHandler) renderTemplate(c *gin.Context, name string, data map[string]interface{}) {
	// Always attach the user's information
	if data == nil {
		data = make(map[string]interface{})
	}

	data["CurrentUser"] = v.getUserFromCookie(c)
	data["Flash"] = c.Query("flash")

	f := template.FuncMap{
		"navbar": func() (template.HTML, error) {
			var buff bytes.Buffer
			err := v.templs.ExecuteTemplate(&buff, "index/navbar", data)
			return template.HTML(buff.String()), err
		},
	}

	if name != "" {
		f["yield"] = func() (template.HTML, error) {
			var buff bytes.Buffer
			err := v.templs.ExecuteTemplate(&buff, name, data)
			return template.HTML(buff.String()), err
		}
	}

	// Clone the layout.
	clonedLayout, _ := v.layout.Clone()
	clonedLayout.Funcs(f)

	err := clonedLayout.Execute(c.Writer, nil)
	if err != nil {
		e := v.templs.ExecuteTemplate(c.Writer, "other/error", map[string]interface{}{
			"Error": err.Error(),
		})
		if e != nil {
			golog.Error(err)
		}
	}
}

func (v *ViewHandler) getUserFromCookie(c *gin.Context) *user.User {
	sess := v.getSessionFromRequest(c)
	if sess == nil {
		return nil
	}

	usr, err := v.usrService.GetByUserID(c.Request.Context(), sess.UserID)
	if err != nil {
		golog.Debug("while getting user by ID:", err)
		return nil
	}
	return usr
}

func (v *ViewHandler) getSessionFromRequest(c *gin.Context) *session.Session {
	cookieVal, _ := c.Cookie(session.CookieName)
	sess, err := v.sessionService.Find(c.Request.Context(), cookieVal)
	if err != nil {
		return nil
	}
	return sess
}

package view

import (
  "bytes"
  "github.com/gin-gonic/gin"
  "github.com/jayvib/golog"
  "golang.org/x/crypto/bcrypt"
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

func RegisterRoutes(r gin.IRouter, userService user.Service, sessionService session.Service, templatesGlob, layoutPath string) {
  h := NewHandler(userService, sessionService, templatesGlob, layoutPath)

  // Asset handler
  r.StaticFS("/assets", http.Dir("assets/"))

  // View handlers
  r.GET("/", h.HomePage)
  r.GET("/signup", h.SignupPage)
  r.GET("/login", h.LoginPage)
  r.GET("/account", h.EditUserPage)
  r.GET("/signout", h.SignOutPage)
  r.GET("/images/new", h.UploadImagePage)

  // Controller handler
  r.POST("/signup", h.HandleSignUp)
  r.POST("/login", h.HandleLogin)
  r.POST("/account", h.HandleEditUser)
  r.POST("/images/new", h.HandleImageUpload)
}

func NewHandler(userService user.Service, sessionService session.Service, templatesGlob, layoutPath string) *ViewHandler {
  return &ViewHandler{
    usrService: userService,
    sessionService: sessionService,
    templs: template.Must(template.ParseGlob(templatesGlob)),
    layout: template.Must(template.New("layout.html").
      Funcs(funcs).ParseFiles(layoutPath)),
  }
}

type ViewHandler struct {
  templs *template.Template
  layout *template.Template
  usrService user.Service
  sessionService session.Service
}

// #################CONTROLLERS################
func (v *ViewHandler) HandleSignUp(c *gin.Context) {
  email := c.PostForm("email")
  username := c.PostForm("username")
  password := c.PostForm("password")
  usr := &user.User{
    Email: email,
    Username: username,
    Password: password,
  }

  err := v.usrService.Register(c.Request.Context(), usr)
  if err != nil {
    v.renderTemplate(c, "users/signup", map[string]interface{}{
      "User": usr,
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
	//usr := v.getUserFromCookie(c)
	//img := &image.Image{ UserID: usr.UserID }

	// Create an image object


}

func (v *ViewHandler) createImageFromFile(c *gin.Context) {
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
  // Get the user detail through username
  err := v.usrService.Login(c.Request.Context(), usr)
  if err != nil {
    v.renderTemplate(c, "sessions/login", map[string]interface{}{
      "Error": getMessage(err),
      "User": usr,
      "Next": next,
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
    Email: email,
    Username: usr.Username,
    Password: currentPassword,
  }

  if newPassword != "" {
    err := v.usrService.Login(c.Request.Context(), tmpUser)
    if err != nil {
      v.renderTemplate(c, "users/edit", map[string]interface{}{
        "Error": getMessage(err),
        "User": tmpUser,
      })
      return
    }

    // ===============NOT SURE IF IT BELONGS HERE==============
    hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
      err := user.NewError(err)
      v.renderTemplate(c, "users/edit", map[string]interface{}{
        "Error": getMessage(err),
        "User": tmpUser,
      })
    }
    usr.Password = string(hash)
    // ========================================================
  }

  usr.Email = email

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
  type messenger interface{
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
  v.renderTemplate(c, "index/navbar", nil)
}

// ###################VIEW####################

func (v *ViewHandler) SignupPage(c *gin.Context) {
  v.renderTemplate(c, "users/signup", nil)
}

func (v *ViewHandler) LoginPage(c *gin.Context) {
  next := c.Query("next")
  v.renderTemplate(c, "sessions/login",  map[string]interface{}{
    "Next": next,
  })
}

func (v *ViewHandler) EditUserPage(c *gin.Context) {
  usr := v.getUserFromCookie(c)
  v.renderTemplate(c, "users/edit", map[string]interface{}{
    "User": usr,
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

  usr, err := v.usrService.GetByID(c.Request.Context(), sess.UserID)
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


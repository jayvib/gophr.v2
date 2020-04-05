package view

import (
  "bytes"
  "github.com/gin-gonic/gin"
  "github.com/jayvib/golog"
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
  r.StaticFS("/assets", http.Dir("assets/"))
  r.GET("/", h.HomePage)
  r.GET("/signup", h.Signup)
  r.GET("/login", h.Login)
  r.POST("/signup", h.HandleSignUp)
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
func (v *ViewHandler) Signup(c *gin.Context) {
  v.renderTemplate(c, "users/signup", nil)
}
func (v *ViewHandler) Login(c *gin.Context) {
  v.renderTemplate(c, "sessions/login", nil)
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
  cookieVal, _ := c.Cookie(session.CookieName)
  sess, err := v.sessionService.Find(c.Request.Context(), cookieVal)
  if err != nil {
    return nil
  }
  usr, err := v.usrService.GetByID(c.Request.Context(), sess.ID)
  if err != nil {
    return nil
  }
  return usr
}


package view

import (
  "bytes"
  "github.com/gin-gonic/gin"
  "github.com/jayvib/golog"
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

func RegisterRoutes(r gin.IRouter, userService user.Service, templatesGlob, layoutPath string) {
  h := NewHandler(userService, templatesGlob, layoutPath)
  r.StaticFS("/assets", http.Dir("assets/"))
  r.GET("/", h.HomePage)
  r.GET("/signup", h.Signup)
  r.GET("/login", h.Login)
  r.POST("/signup", h.HandleSignUp)
}

func NewHandler(userService user.Service, templatesGlob, layoutPath string) *ViewHandler {
  return &ViewHandler{
    usrService: userService,
    templs: template.Must(template.ParseGlob(templatesGlob)),
    layout: template.Must(template.New("layout.html").
      Funcs(funcs).ParseFiles(layoutPath)),
  }
}

type ViewHandler struct {
  templs *template.Template
  layout *template.Template
  usrService user.Service
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
    v.renderTemplate(c, "other/error", map[string]interface{}{
      "Error": getMessage(err),
    })
    return
  }
  c.Redirect(http.StatusFound, "/?flash=User+created")
}

func getMessage(err error) string {
  var message string
  if uerr, ok := err.(*user.Error); ok {
    message = uerr.Message()
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


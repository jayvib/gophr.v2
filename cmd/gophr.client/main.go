package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
)

var funcs = template.FuncMap{
	"yield": func() (string, error) {
		return "", nil
	},
	"navbar": func() (string, error) {
		return "", nil
	},
	"pagename": func() string {
		return ""
	},
}

func RegisterRoutes(r gin.IRouter, templatesGlob, layoutPath string) {
	gophr := NewGophr(templatesGlob, layoutPath)

	r.StaticFS("/assets", http.Dir("v2/assets/"))

	r.GET("/", gophr.HomePage)
	r.GET("/login", gophr.LoginPage)
}

func NewGophr(templateGlob, layoutPath string) *Gophr {
	return &Gophr{
		tmpls: template.Must(template.ParseGlob(templateGlob)),
		layout: template.Must(template.New("layout.html").Funcs(funcs).ParseFiles(layoutPath)),
	}
}

type Gophr struct {
	tmpls *template.Template
	layout *template.Template
}

func (g *Gophr) HomePage(c *gin.Context) {
	g.renderTeamplate(c, "index/home", "Gophr", nil)
}

func (g *Gophr) LoginPage(c *gin.Context) {
	g.renderTeamplate(c, "session/login", "Login", nil)
}

func (g *Gophr) renderTeamplate(c *gin.Context, name string, pageName string, data map[string]interface{}) {
	f := template.FuncMap{
		"navbar": func() (template.HTML, error) {
			var buff bytes.Buffer
			err := g.tmpls.ExecuteTemplate(&buff, "index/navbar", data)
			return template.HTML(buff.String()), err
		},
	}

	if name != "" {
		f["yield"] = func() (template.HTML, error) {
			var buff bytes.Buffer
			err := g.tmpls.ExecuteTemplate(&buff, name, data)
			return template.HTML(buff.String()), err
		}
	}

	if pageName != "" {
		f["pagename"] = func() string {
			return pageName
		}
	}

	clonedLayout, _ := g.layout.Clone()
	clonedLayout.Funcs(f)

	err := clonedLayout.Execute(c.Writer, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	engine := gin.Default()
	RegisterRoutes(engine, "v2/templates/**/*.html", "v2/templates/layout.html")
	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

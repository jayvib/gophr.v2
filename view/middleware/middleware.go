package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	"gophr.v2/session"
	"net/http"
	"net/url"
	"strings"
)

func RequireLogin(sessionService session.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the cookie
		golog.Debug(c.Request.Cookies())
		cookieValue, err := c.Cookie(session.CookieName)
		golog.Debug(cookieValue, err)
		if err != nil {
			redirectToLogin(c)
			return
		}
		golog.Debugf("%#v", cookieValue)
		// Get the detail of the session
		sess, err := sessionService.Find(c.Request.Context(), cookieValue)
		if err != nil {
			redirectToLogin(c)
			return
		}
		golog.Debugf("%#v", sess)
		if sess.IsExpired() {
			redirectToLogin(c)
			return
		}
		c.Next()
	}
}

func redirectToLogin(c *gin.Context) {
	next := url.Values{}
	trimmedUrl := strings.TrimLeftFunc(c.Request.URL.String(), func(r rune) bool {
		if r == '/' {
			return true
		}
		return false
	})
	next.Add("next", url.QueryEscape(trimmedUrl))
	c.Redirect(http.StatusTemporaryRedirect, "/login?"+next.Encode())
	c.Abort()
}

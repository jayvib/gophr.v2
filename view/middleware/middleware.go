package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	"gophr.v2/session"
	"net/http"
	"net/url"
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
	next.Add("next", url.QueryEscape(c.Request.URL.String()))
	c.Redirect(http.StatusTemporaryRedirect, "/login?"+next.Encode())
	c.Abort()
}

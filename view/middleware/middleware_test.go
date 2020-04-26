//+build unit

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophr.v2/http/httputil"
	"gophr.v2/session"
	"gophr.v2/session/mocks"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func TestRequireLogin(t *testing.T) {

	t.Run("Active Session", func(t *testing.T) {
		sess := &session.Session{
			ID:     "testing123",
			Expiry: time.Now().Add(session.Duration),
		}
		svc := new(mocks.Service)
		svc.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(sess, nil).Once()
		r := gin.New()
		authRouter := r.Use(RequireLogin(svc))
		authRouter.GET("/hello", func(c *gin.Context) { c.Status(http.StatusOK) })
		resp := httputil.PerformRequest(r, http.MethodGet, "/hello", nil, func(r *http.Request) {
			cookie := &http.Cookie{
				Name:    session.CookieName,
				Value:   sess.ID,
				Expires: sess.Expiry,
			}
			r.AddCookie(cookie)
		})
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Expired Session", func(t *testing.T) {
		sess := &session.Session{
			ID:     "testing123",
			Expiry: time.Now().AddDate(0, 0, -1),
		}
		svc := new(mocks.Service)
		svc.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(sess, nil).Once()
		r := gin.New()
		authRouter := r.Use(RequireLogin(svc))
		authRouter.GET("/hello", func(c *gin.Context) { c.Status(http.StatusOK) })
		resp := httputil.PerformRequest(r, http.MethodGet, "/hello", nil, func(r *http.Request) {
			cookie := &http.Cookie{
				Name:    session.CookieName,
				Value:   sess.ID,
				Expires: sess.Expiry,
			}
			r.AddCookie(cookie)
		})
		assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
	})

	t.Run("Cookie doesn't exists", func(t *testing.T) {
		sess := &session.Session{
			ID:     "testing123",
			Expiry: time.Now(),
		}
		svc := new(mocks.Service)
		svc.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(sess, nil).Once()
		r := gin.New()
		authRouter := r.Use(RequireLogin(svc))
		authRouter.GET("/hello", func(c *gin.Context) { c.Status(http.StatusOK) })
		resp := httputil.PerformRequest(r, http.MethodGet, "/hello", nil)
		assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
	})

	t.Run("Session doesn't exists", func(t *testing.T) {
		sess := &session.Session{
			ID:     "testing123",
			Expiry: time.Now(),
		}
		svc := new(mocks.Service)
		svc.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(nil, session.ErrNotFound).Once()
		r := gin.New()
		authRouter := r.Use(RequireLogin(svc))
		authRouter.GET("/hello", func(c *gin.Context) { c.Status(http.StatusOK) })
		resp := httputil.PerformRequest(r, http.MethodGet, "/hello", nil, func(r *http.Request) {
			cookie := &http.Cookie{
				Name:    session.CookieName,
				Value:   sess.ID,
				Expires: sess.Expiry,
			}
			r.AddCookie(cookie)
		})
		assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
	})

}

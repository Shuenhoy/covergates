package web

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/drone/go-login/login"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	keyLogin    = "login"
	keyAccess   = "access"
	keyRefresh  = "refresh"
	keyExpires  = "expires"
	keyRedirect = "redirect"
)

// HandleLogin user
func HandleLogin(
	config *config.Config,
	scm core.SCMProvider,
	scmService core.SCMService,
	session core.Session,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool(keyLogin) {
			return
		}
		ctx := c.Request.Context()
		client, err := scmService.Client(scm)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		token := TokenFrom(c)
		user := session.GetUser(c)
		if session.ShouldBindUser(c) {
			user, err = client.Users().Bind(ctx, user, token)
			session.EndBindUser(c)
		} else {
			user, err = createOrUpdateUser(ctx, client, token)
		}

		if err != nil {
			log.Error(err)
			c.String(400, err.Error())
		}
		if err := session.CreateUser(c, user); err != nil {
			log.Error(err)
			c.String(400, err.Error())
			return
		}

		redirect := c.GetString(keyRedirect)
		if redirect == "" {
			redirect = "/"
		}

		c.Redirect(301, strings.TrimRight(config.Server.BaseURL(), "/")+redirect)

	}
}

func createOrUpdateUser(ctx context.Context, client core.Client, token *core.Token) (*core.User, error) {
	user, err := client.Users().Find(ctx, token)
	if err != nil {
		user, err = client.Users().Create(ctx, token)
	} else {
		user, err = client.Users().Update(ctx, token)
	}
	return user, err
}

// MiddlewareLogin context
func MiddlewareLogin(scm core.SCMProvider, m core.LoginMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		// if query redirect exsit, write it to cookie
		redirect := c.Query("redirect")
		if redirect != "" {
			redirect, err := url.QueryUnescape(redirect)
			if err == nil {
				http.SetCookie(c.Writer, &http.Cookie{
					Name:   keyRedirect,
					Value:  redirect,
					MaxAge: 1800,
				})
			}
		}
		middleware := m.Handler(scm)
		h := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			err := login.ErrorFrom(ctx)
			if err != nil {
				c.Error(err)
				c.Abort()
				return
			}
			tok := login.TokenFrom(ctx)
			c.Set(keyLogin, true)
			c.Set(keyAccess, tok.Access)
			c.Set(keyExpires, tok.Expires)
			c.Set(keyRefresh, tok.Refresh)

			cookie, err := r.Cookie(keyRedirect)
			if err == nil {
				c.Set(keyRedirect, cookie.Value)
				http.SetCookie(w, &http.Cookie{
					Name:    keyRedirect,
					MaxAge:  -1,
					Expires: time.Unix(0, 0),
				})
			}
		}))
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// MiddlewareBindUser handle bind user request
func MiddlewareBindUser(session core.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, bind := c.GetQuery("bind")
		if !bind {
			return
		}
		user := session.GetUser(c)
		if user.Login == "" {
			c.String(401, "Unauthorized")
			c.Abort()
			return
		}
		session.StartBindUser(c)
	}
}

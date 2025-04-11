package middleware

import (
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/pkg/rest/serverinterceptor"

	"net/http"
)

type UserLoginMiddleware struct {
	c config.Config
}

func NewUserLoginMiddleware(c config.Config) *UserLoginMiddleware {
	return &UserLoginMiddleware{
		c: c,
	}
}

func (m *UserLoginMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	authorizeHandle := serverinterceptor.AuthorizeHandle(m.c.JwtSecret, true)
	return func(w http.ResponseWriter, r *http.Request) {
		authorizeHandle(w, r, next)
	}
}

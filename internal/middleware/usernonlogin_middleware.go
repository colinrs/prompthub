package middleware

import (
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/pkg/rest/serverinterceptor"

	"net/http"
)

type UserNonLoginMiddleware struct {
	c config.Config
}

func NewUserNonLoginMiddleware(c config.Config) *UserNonLoginMiddleware {
	return &UserNonLoginMiddleware{
		c: c,
	}
}

func (m *UserNonLoginMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	authorizeHandle := serverinterceptor.AuthorizeHandle(m.c.JwtSecret, false)
	return func(w http.ResponseWriter, r *http.Request) {
		authorizeHandle(w, r, next)
	}
}

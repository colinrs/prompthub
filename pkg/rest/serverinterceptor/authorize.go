package serverinterceptor

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"

	"github.com/colinrs/prompthub/pkg/constant"
	"github.com/colinrs/prompthub/pkg/response"
	"github.com/colinrs/prompthub/pkg/utils"
	"github.com/zeromicro/go-zero/core/logc"
)

const (
	noDetailReason = "no detail reason"
)

var (
	errInvalidToken = errors.New("invalid auth token")
)

type (
	AuthorizeOptions struct {
		Secret []byte
	}
	// UnauthorizedCallback defines the method of unauthorized callback.
	UnauthorizedCallback func(w http.ResponseWriter, r *http.Request, err error)
	// AuthorizeOption defines the method to customize an AuthorizeOptions.
	AuthorizeOption func(opts *AuthorizeOptions)
)

// Authorize returns an authorization middleware.
func Authorize(secret string, forceLogin bool, opts ...AuthorizeOption) func(http.Handler) http.Handler {
	var authOpts AuthorizeOptions
	authOpts.Secret = []byte(secret)
	for _, opt := range opts {
		opt(&authOpts)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if len(token) == 0 {
				next.ServeHTTP(w, r)
				return
			}
			Bearer := "Bearer "
			if len(token) < len(Bearer) || token[:len(Bearer)] != Bearer && forceLogin {
				unauthorized(w, r, errInvalidToken, nil)
				return
			}
			token = token[len(Bearer):]
			claims, ok := utils.ParseJWT(token, authOpts.Secret)
			if !ok && forceLogin {
				unauthorized(w, r, errInvalidToken, nil)
				return
			}
			ctx := r.Context()
			for k, v := range claims {
				switch k {
				case constant.UserId, constant.UserName, constant.Email:
					ctx = context.WithValue(ctx, k, v)
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthorizeHandle(secret string, forceLogin bool, opts ...AuthorizeOption) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var authOpts AuthorizeOptions
	authOpts.Secret = []byte(secret)
	for _, opt := range opts {
		opt(&authOpts)
	}
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			next.ServeHTTP(w, r)
			return
		}
		Bearer := "Bearer "
		if len(token) < len(Bearer) || token[:len(Bearer)] != Bearer && forceLogin {
			unauthorized(w, r, errInvalidToken, nil)
			return
		}
		token = token[len(Bearer):]
		claims, ok := utils.ParseJWT(token, authOpts.Secret)
		if !ok && forceLogin {
			unauthorized(w, r, errInvalidToken, nil)
			return
		}
		ctx := r.Context()
		for k, v := range claims {
			switch k {
			case constant.UserId, constant.UserName, constant.Email:
				ctx = context.WithValue(ctx, k, v)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// WithSecret returns an AuthorizeOption with setting previous secret.
func WithSecret(secret string) AuthorizeOption {
	return func(opts *AuthorizeOptions) {
		opts.Secret = []byte(secret)
	}
}

func detailAuthLog(r *http.Request, reason string) {
	// discard dump error, only for debug purpose
	details, _ := httputil.DumpRequest(r, true)
	logc.Errorf(r.Context(), "authorize failed: %s\n=> %+v", reason, string(details))
}

func unauthorized(w http.ResponseWriter, r *http.Request, err error, callback UnauthorizedCallback) {
	writer := response.NewHeaderOnceResponseWriter(w)

	if err != nil {
		detailAuthLog(r, err.Error())
	} else {
		detailAuthLog(r, noDetailReason)
	}
	// let callback go first, to make sure we respond with user-defined HTTP header
	if callback != nil {
		callback(writer, r, err)
	}
	// if user not setting HTTP header, we set header with 401
	writer.WriteHeader(http.StatusUnauthorized)
}

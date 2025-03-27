package user

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/user"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewRefreshTokenLogic(r.Context(), svcCtx)
		resp, err := l.RefreshToken()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

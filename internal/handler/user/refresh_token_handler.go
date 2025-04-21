package user

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/user"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/httpy"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshTokenRequeste
		if err := httpx.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		l := user.NewRefreshTokenLogic(r.Context(), svcCtx)
		resp, err := l.RefreshToken(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}

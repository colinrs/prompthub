package user

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/user"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/httpy"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func VerificationCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VerificationCodeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := user.NewVerificationCodeLogic(r.Context(), svcCtx)
		resp, err := l.VerificationCode(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}

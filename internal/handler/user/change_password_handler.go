package user

import (
	"github.com/colinrs/prompthub/pkg/httpy"
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/user"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangePasswordRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		l := user.NewChangePasswordLogic(r.Context(), svcCtx)
		err := l.ChangePassword(&req)
		httpy.ResultCtx(r, w, nil, err)
	}
}

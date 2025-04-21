package user

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/user"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/httpy"
)

func RegisterUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterUserRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		if svcCtx.DetectorSWD.Detect(req.Name + req.Email) {
			httpy.ResultCtx(r, w, nil, code.ErrSensitiveWord)
			return
		}
		l := user.NewRegisterUserLogic(r.Context(), svcCtx)
		resp, err := l.RegisterUser(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}

package user

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/user"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/pkg/httpy"
)

func GetUserSavePromptsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserSavePromptsLogic(r.Context(), svcCtx)
		resp, err := l.GetUserSavePrompts()
		httpy.ResultCtx(r, w, resp, err)
	}
}

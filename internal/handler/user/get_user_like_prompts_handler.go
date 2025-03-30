package user

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/user"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/pkg/httpy"
)

func GetUserLikePromptsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserLikePromptsLogic(r.Context(), svcCtx)
		resp, err := l.GetUserLikePrompts()
		httpy.ResultCtx(r, w, resp, err)
	}
}

package user

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/user"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/pkg/httpy"
)

func GetUserPromptsListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserPromptsListLogic(r.Context(), svcCtx)
		resp, err := l.GetUserPromptsList()
		httpy.ResultCtx(r, w, resp, err)
	}
}

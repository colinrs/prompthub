package prompt

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/prompt"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreatePromptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatePromptRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := prompt.NewCreatePromptLogic(r.Context(), svcCtx)
		err := l.CreatePrompt(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}

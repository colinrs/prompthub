package prompt

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/prompt"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/httpy"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SavePromptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SavePromptRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		l := prompt.NewSavePromptLogic(r.Context(), svcCtx)
		resp, err := l.SavePrompt(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}

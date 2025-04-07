package prompt

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/prompt"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/httpy"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreatePromptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatePromptRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		l := prompt.NewCreatePromptLogic(r.Context(), svcCtx)
		err := l.CreatePrompt(&req)
		httpy.ResultCtx(r, w, nil, err)
	}
}

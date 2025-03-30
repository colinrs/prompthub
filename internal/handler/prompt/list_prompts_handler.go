package prompt

import (
	"github.com/colinrs/prompthub/pkg/httpy"
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/prompt"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListPromptsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListPromptsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		l := prompt.NewListPromptsLogic(r.Context(), svcCtx)
		resp, err := l.ListPrompts(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}

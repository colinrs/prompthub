package prompt

import (
	"github.com/colinrs/prompthub/pkg/code"
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/prompt"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/httpy"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdatePromptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdatePromptRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		if svcCtx.DetectorSWD.Detect(req.Title + req.Content) {
			httpy.ResultCtx(r, w, nil, code.ErrSensitiveWord)
			return
		}
		l := prompt.NewUpdatePromptLogic(r.Context(), svcCtx)
		err := l.UpdatePrompt(&req)
		httpy.ResultCtx(r, w, nil, err)
	}
}

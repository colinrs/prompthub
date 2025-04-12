package category

import (
	"github.com/colinrs/prompthub/pkg/code"
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/category"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/httpy"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateCategoryRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		if svcCtx.DetectorSWD.Detect(req.Name + req.Color) {
			httpy.ResultCtx(r, w, nil, code.ErrSensitiveWord)
		}
		l := category.NewCreateCategoryLogic(r.Context(), svcCtx)
		err := l.CreateCategory(&req)
		httpy.ResultCtx(r, w, nil, err)
	}
}

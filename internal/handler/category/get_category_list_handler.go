package category

import (
	"net/http"

	"github.com/colinrs/prompthub/internal/logic/category"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/httpy"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCategoryListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCategoryListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		l := category.NewGetCategoryListLogic(r.Context(), svcCtx)
		resp, err := l.GetCategoryList(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}

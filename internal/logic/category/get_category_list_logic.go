package category

import (
	"context"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryListLogic {
	return &GetCategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryListLogic) GetCategoryList(req *types.GetCategoryListRequest) (resp *types.CategoryResponse, err error) {
	categoryTable := gen.Use(l.svcCtx.DB).CategoryTable
	offset, limit := utils.Page(req.Page, req.PageSize)
	categoryList, err := categoryTable.Where(categoryTable.CategoryStatus.Eq(config.CategoryStatusNormal)).
		Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, err
	}
	total, err := categoryTable.Where(categoryTable.CategoryStatus.Eq(config.CategoryStatusNormal)).Count()
	if err != nil {
		return nil, err
	}
	resp = &types.CategoryResponse{
		List:     []types.Category{},
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    int(total),
	}
	for _, category := range categoryList {
		resp.List = append(resp.List, types.Category{
			Id:    int64(category.ID),
			Name:  category.CategoryName,
			Color: category.Color,
		})
	}
	return
}

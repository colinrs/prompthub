package category

import (
	"context"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/model"
	"github.com/colinrs/prompthub/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryRequest) error {
	categoryTable := gen.Use(l.svcCtx.DB).CategoryTable
	categories, err := categoryTable.Where(categoryTable.CategoryStatus.Eq(config.CategoryStatusNormal),
		categoryTable.CategoryName.Eq(req.Name)).Find()
	if err != nil {
		return err
	}
	if len(categories) > 0 {
		return code.ErrCategoryAlreadyExists
	}
	category := &model.CategoryTable{
		CategoryName:   req.Name,
		Color:          req.Color,
		CategoryStatus: config.CategoryStatusNormal,
	}
	if err := categoryTable.Create(category); err != nil {
		return err
	}
	return nil
}

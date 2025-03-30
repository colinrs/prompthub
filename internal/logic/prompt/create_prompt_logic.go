package prompt

import (
	"context"
	"errors"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/model"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreatePromptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePromptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePromptLogic {
	return &CreatePromptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePromptLogic) CreatePrompt(req *types.CreatePromptRequest) error {
	userID, err := utils.GetUserIDFromCtx(l.ctx)
	if err != nil {
		return err
	}
	categoryTable := gen.Use(l.svcCtx.DB).CategoryTable
	category, err := categoryTable.Where(categoryTable.ID.Eq(int32(req.CategoryId))).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if category == nil {
		return code.ErrParam
	}
	promptsTable := gen.Use(l.svcCtx.DB).PromptsTable
	prompts := &model.PromptsTable{
		ID:            0,
		Title:         req.Title,
		Content:       req.Title,
		Category:      int32(req.CategoryId),
		PromptsStatus: config.PromptsStatusNormal,
		CreatedBy:     int32(userID),
	}
	err = promptsTable.Create(prompts)
	if err != nil {
		return err
	}
	return nil
}

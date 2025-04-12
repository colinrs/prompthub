package prompt

import (
	"context"
	"errors"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/pkg/code"
	"gorm.io/gorm"

	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePromptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePromptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePromptLogic {
	return &UpdatePromptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePromptLogic) UpdatePrompt(req *types.UpdatePromptRequest) error {
	categoryTable := gen.Use(l.svcCtx.DB).CategoryTable
	category, err := categoryTable.Where(categoryTable.ID.Eq(int32(req.CategoryId))).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if category == nil {
		return code.ErrParam
	}

	promptsTable := gen.Use(l.svcCtx.DB).PromptsTable
	prompt, err := promptsTable.Where(promptsTable.ID.Eq(int32(req.ID))).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if prompt == nil {
		return code.ErrParam
	}
	prompt.Title = req.Title
	prompt.Content = req.Content
	prompt.Category = int32(req.CategoryId)
	err = promptsTable.Where(promptsTable.ID.Eq(int32(req.ID))).Save(prompt)
	if err != nil {
		return err
	}

	return nil
}

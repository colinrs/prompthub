package prompt

import (
	"context"
	"github.com/colinrs/prompthub/gen"

	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePromptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePromptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePromptLogic {
	return &DeletePromptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePromptLogic) DeletePrompt(req *types.DeletePromptRequest) error {
	promptsTable := gen.Use(l.svcCtx.DB).PromptsTable
	_, err := promptsTable.Where(promptsTable.ID.Eq(int32(req.PromptID))).Delete()
	return err
}

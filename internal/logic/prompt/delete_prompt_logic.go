package prompt

import (
	"context"

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
	// todo: add your logic here and delete this line

	return nil
}

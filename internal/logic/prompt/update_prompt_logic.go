package prompt

import (
	"context"

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
	// todo: add your logic here and delete this line

	return nil
}

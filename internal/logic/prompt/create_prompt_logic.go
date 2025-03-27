package prompt

import (
	"context"

	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return nil
}

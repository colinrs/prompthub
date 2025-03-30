package prompt

import (
	"context"

	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPromptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPromptLogic {
	return &GetPromptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPromptLogic) GetPrompt(req *types.GetPromptRequest) (resp *types.Prompt, err error) {
	// todo: add your logic here and delete this line

	return
}

package prompt

import (
	"context"
	"errors"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/model"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/utils"
	"gorm.io/gorm"

	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SavePromptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSavePromptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePromptLogic {
	return &SavePromptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SavePromptLogic) SavePrompt(req *types.SavePromptRequest) (resp *types.SavePromptResponse, err error) {
	userId, err := utils.GetUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return nil, code.ErrUserNoLogin
	}
	promptsTable := gen.Use(l.svcCtx.DB).PromptsTable
	prompt, err := promptsTable.Where(promptsTable.PromptsStatus.Eq(config.PromptsStatusNormal),
		promptsTable.ID.Eq(int32(req.PromptID))).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || prompt == nil {
		return nil, code.PromptNotFound
	}
	saveTable := gen.Use(l.svcCtx.DB).UsersSave
	save := model.UsersSave{
		UserID:    int32(userId),
		PromptsID: int32(req.PromptID),
	}
	if req.Action == config.UnSaveAction {
		_, err = saveTable.Where(
			saveTable.UserID.Eq(int32(userId)), saveTable.PromptsID.Eq(int32(req.PromptID))).Unscoped().Delete()
		if err != nil {
			return nil, err
		}
		return
	}
	err = saveTable.Save(&save)
	if err != nil {
		return nil, err
	}
	return
}

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

type LikePromptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikePromptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikePromptLogic {
	return &LikePromptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikePromptLogic) LikePrompt(req *types.LikePromptRequest) (resp *types.LikePromptResponse, err error) {
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
	likeTable := gen.Use(l.svcCtx.DB).UsersLike
	like := model.UsersLike{
		UserID:    int32(userId),
		PromptsID: int32(req.PromptID),
	}
	countTable := gen.Use(l.svcCtx.DB).PromptsCountTable

	if req.Action == config.LikeAction {
		err = likeTable.Save(&like)
		if err != nil {
			return nil, err
		}
		_ = countTable.Create(&model.PromptsCountTable{
			PromptsID: int32(req.PromptID),
		})
		_, err = countTable.Where(countTable.PromptsID.Eq(int32(req.PromptID))).Update(countTable.LikeCount, countTable.LikeCount.Add(1))
		if err != nil {
			return nil, err
		}
	}
	if req.Action == config.UnlikeAction {
		_, err = likeTable.Where(
			likeTable.UserID.Eq(int32(userId)), likeTable.PromptsID.Eq(int32(req.PromptID))).Unscoped().Delete()
		if err != nil {
			return nil, err
		}
		_ = countTable.Create(&model.PromptsCountTable{
			PromptsID: int32(req.PromptID),
		})
		_, err = countTable.Where(countTable.PromptsID.Eq(int32(req.PromptID)), countTable.LikeCount.Gt(0)).Update(countTable.LikeCount, countTable.LikeCount.Add(-1))
		if err != nil {
			return nil, err
		}
	}

	return
}

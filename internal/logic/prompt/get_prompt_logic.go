package prompt

import (
	"context"
	"github.com/colinrs/prompthub/internal/manager"
	"github.com/colinrs/prompthub/pkg/code"
	"golang.org/x/sync/errgroup"

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
	promoHabit := make([]*manager.PromptHabit, 0, 1)
	promoHabit = append(promoHabit, &manager.PromptHabit{
		PromptID: req.PromptID,
	})
	g := errgroup.Group{}
	g.Go(func() error {
		manager.PromptListDetail(promoHabit).SetSave(l.ctx, l.svcCtx)
		return nil
	})
	g.Go(func() error {
		manager.PromptListDetail(promoHabit).SetLikes(l.ctx, l.svcCtx)
		return nil
	})
	g.Go(func() error {
		manager.PromptListDetail(promoHabit).SetCount(l.ctx, l.svcCtx)
		return nil
	})
	_ = g.Wait()
	promptManger := manager.NewPrompt(l.ctx, l.svcCtx)
	promptsDetail, err := promptManger.GetPromptDetailByIds([]int32{int32(req.PromptID)})
	if err != nil {
		return nil, err
	}
	if len(promptsDetail) != 1 {
		return nil, code.PromptNotFound
	}
	item := promptsDetail[0]
	liked := false
	saved := false
	likes := 0
	if habit, ok := manager.PromptListDetail(promoHabit).Map()[int32(item.PromptID)]; ok {
		liked = habit.Liked
		saved = habit.Saved
		likes = int(habit.LikeCount)
	}
	resp = &types.Prompt{
		Id:            item.PromptID,
		Title:         item.Title,
		Content:       item.Content,
		CategoryID:    item.CategoryID,
		CategoryColor: item.CategoryColor,
		Category:      item.Category,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		CreatedBy: types.CreatedBy{
			UserID:   item.CreatedBy.UserID,
			UserName: item.CreatedBy.UserName,
			Avatar:   item.CreatedBy.Avatar,
		},
		Likes: likes,
		Liked: liked,
		Saved: saved,
	}
	return
}

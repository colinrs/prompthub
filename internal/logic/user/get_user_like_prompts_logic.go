package user

import (
	"context"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/manager"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/utils"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
)

type GetUserLikePromptsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLikePromptsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLikePromptsLogic {
	return &GetUserLikePromptsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLikePromptsLogic) GetUserLikePrompts() (resp *types.UserPromptsResponse, err error) {
	usersLikeTable := gen.Use(l.svcCtx.DB).UsersLike
	userID, err := utils.GetUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	usersLike, err := usersLikeTable.Where(usersLikeTable.UserID.Eq(int32(userID))).Find()
	if err != nil {
		return nil, err
	}
	promptIds := make([]int32, 0, len(usersLike))
	for _, l := range usersLike {
		promptIds = append(promptIds, l.PromptsID)
	}
	prompt := manager.NewPrompt(l.ctx, l.svcCtx)
	promptsDetail, err := prompt.GetPromptDetailByIds(promptIds)
	if err != nil {
		return nil, err
	}
	promoHabit := make([]*manager.PromptHabit, 0, len(promptsDetail))
	for _, l := range promptsDetail {
		promoHabit = append(promoHabit, &manager.PromptHabit{
			PromptID: l.PromptID,
		})
	}
	g := errgroup.Group{}
	g.Go(func() error {
		manager.PromptListDetail(promoHabit).SetSave(l.ctx, l.svcCtx)
		return nil
	})
	g.Go(func() error {
		manager.PromptListDetail(promoHabit).SetCount(l.ctx, l.svcCtx)
		return nil
	})
	_ = g.Wait()
	resp = &types.UserPromptsResponse{
		List: lo.Map(promptsDetail, func(item *manager.PromptDetail, index int) types.Prompt {
			saved := false
			likes := 0
			if habit, ok := manager.PromptListDetail(promoHabit).Map()[int32(item.PromptID)]; ok {
				saved = habit.Saved
				likes = int(habit.LikeCount)
			}
			return types.Prompt{
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
				Liked: true,
				Saved: saved,
			}
		}),
		Page:     1,
		PageSize: 10,
		Total:    len(promptsDetail),
	}
	return
}

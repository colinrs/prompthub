package user

import (
	"context"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/manager"
	"github.com/colinrs/prompthub/model"
	"github.com/colinrs/prompthub/pkg/utils"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPromptsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPromptsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPromptsListLogic {
	return &GetUserPromptsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPromptsListLogic) GetUserPromptsList() (resp *types.UserPromptsResponse, err error) {
	promptsTable := gen.Use(l.svcCtx.DB).PromptsTable
	userID, err := utils.GetUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	userPrompts, err := promptsTable.Where(promptsTable.CreatedBy.Eq(int32(userID))).Find()
	if err != nil {
		return nil, err
	}
	categoryIds := make([]int32, 0, len(userPrompts))
	for _, l := range userPrompts {
		categoryIds = append(categoryIds, l.Category)
	}
	prompt := manager.NewPrompt(l.ctx, l.svcCtx)
	categoryMap, err := prompt.GetCategoryMap(categoryIds)
	if err != nil {
		return nil, err
	}
	promoHabit := make([]*manager.PromptHabit, 0, len(userPrompts))
	for _, l := range userPrompts {
		promoHabit = append(promoHabit, &manager.PromptHabit{
			PromptID: uint(l.ID),
		})
	}
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

	resp = &types.UserPromptsResponse{
		List: lo.Map(userPrompts, func(item *model.PromptsTable, index int) types.Prompt {
			categoryID := 0
			categoryColor := ""
			categoryName := ""
			if category, ok := categoryMap[item.Category]; ok {
				categoryColor = category.Color
				categoryName = category.CategoryName
				categoryID = int(item.Category)
			}
			liked := false
			saved := false
			likes := 0
			if habit, ok := manager.PromptListDetail(promoHabit).Map()[item.ID]; ok {
				liked = habit.Liked
				saved = habit.Saved
				likes = int(habit.LikeCount)
			}
			return types.Prompt{
				Id:            uint(item.ID),
				Title:         item.Title,
				Content:       item.Content,
				CategoryID:    uint(categoryID),
				CategoryColor: categoryColor,
				Category:      categoryName,
				CreatedAt:     item.CreatedAt.UTC().Format(utils.TimeLayout),
				UpdatedAt:     item.UpdatedAt.UTC().Format(utils.TimeLayout),
				CreatedBy: types.CreatedBy{
					UserID:   userID,
					UserName: utils.GetUserNameFromCtx(l.ctx),
					Avatar:   "",
				},
				Likes: likes,
				Liked: liked,
				Saved: saved,
			}
		}),
		Page:     1,
		PageSize: 10,
		Total:    len(userPrompts),
	}
	return
}

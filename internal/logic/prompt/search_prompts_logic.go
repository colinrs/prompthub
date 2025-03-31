package prompt

import (
	"context"
	"golang.org/x/sync/errgroup"

	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/manager"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	gormgen "gorm.io/gen"
)

type SearchPromptsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchPromptsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchPromptsLogic {
	return &SearchPromptsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchPromptsLogic) SearchPrompts(req *types.SearchPromptsRequest) (resp *types.SearchPromptsResponse, err error) {
	promptsTable := gen.Use(l.svcCtx.DB).PromptsTable
	query := promptsTable.Where(promptsTable.PromptsStatus.Eq(config.PromptsStatusNormal))

	likeConds := make([]gormgen.Condition, 0, 2)
	if req.Title != "" {
		likeConds = append(likeConds, promptsTable.Title.Like("%"+req.Title+"%"))
	}
	if req.Content != "" {
		likeConds = append(likeConds, promptsTable.Content.Like("%"+req.Content+"%"))
	}
	if len(likeConds) > 0 {
		query = query.Or(likeConds...)
	}
	if req.CategoryID != 0 {
		query = query.Where(promptsTable.Category.Eq(int32(req.CategoryID)))
	}
	promptList, err := query.Find()
	if err != nil {
		return nil, err
	}
	total, err := query.Count()
	if err != nil {
		return nil, err
	}
	promptIds := make([]int32, 0, len(promptList))
	for _, prompt := range promptList {
		promptIds = append(promptIds, prompt.ID)
	}
	promoHabit := make([]*manager.PromptHabit, 0, len(promptList))
	for _, l := range promptList {
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

	promptManger := manager.NewPrompt(l.ctx, l.svcCtx)
	promptsDetail, err := promptManger.GetPromptDetailByIds(promptIds)
	if err != nil {
	}
	resp = &types.SearchPromptsResponse{
		Total: int(total),
		List: lo.UniqMap(promptsDetail, func(item *manager.PromptDetail, _ int) types.Prompt {
			liked := false
			saved := false
			likes := 0
			if habit, ok := manager.PromptListDetail(promoHabit).Map()[int32(item.PromptID)]; ok {
				liked = habit.Liked
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
				Liked: liked,
				Saved: saved,
			}
		}),
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	return
}

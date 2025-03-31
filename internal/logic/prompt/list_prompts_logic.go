package prompt

import (
	"context"
	"github.com/colinrs/prompthub/internal/manager"
	"github.com/colinrs/prompthub/model"
	"golang.org/x/sync/errgroup"

	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/utils"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListPromptsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPromptsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPromptsLogic {
	return &ListPromptsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPromptsLogic) ListPrompts(req *types.ListPromptsRequest) (resp *types.ListPromptResponse, err error) {
	promptsTable := gen.Use(l.svcCtx.DB).PromptsTable
	offset, limit := utils.Page(req.Page, req.PageSize)
	promptList, err := promptsTable.Where(promptsTable.PromptsStatus.Eq(config.PromptsStatusNormal)).
		Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, err
	}
	total, err := promptsTable.Where(promptsTable.PromptsStatus.Eq(config.PromptsStatusNormal)).Count()
	if err != nil {
		return nil, err
	}
	categoryTable := gen.Use(l.svcCtx.DB).CategoryTable
	categoryIds := make([]int32, 0, len(promptList))
	userIds := make([]int32, 0, len(promptList))
	for _, category := range promptList {
		categoryIds = append(categoryIds, category.Category)
		userIds = append(userIds, category.CreatedBy)
	}
	categoryIds = lo.Uniq(categoryIds)
	categoryList, err := categoryTable.Where(categoryTable.ID.In(categoryIds...)).Find()
	if err != nil {
		return nil, err
	}
	categoryMap := make(map[int32]string, len(categoryList))
	lo.UniqMap(categoryList, func(category *model.CategoryTable, _ int) string {
		categoryMap[category.ID] = category.CategoryName
		return ""
	})
	userIds = lo.Uniq(userIds)
	userTable := gen.Use(l.svcCtx.DB).UsersTable
	userList, err := userTable.Where(userTable.UserStatus.Eq(config.UserStatusNormal), userTable.ID.In(userIds...)).Find()
	if err != nil {
		return nil, err
	}
	userMap := make(map[int32]*model.UsersTable, len(userList))
	lo.UniqMap(userList, func(user *model.UsersTable, _ int) string {
		userMap[user.ID] = user
		return ""
	})
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

	resp = &types.ListPromptResponse{
		List:     []types.Prompt{},
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    int(total),
	}
	for _, prompt := range promptList {
		user, ok := userMap[prompt.CreatedBy]
		userName := ""
		avatar := ""
		if ok {
			userName = user.UserName
			avatar = user.Avatar
		}
		liked := false
		saved := false
		likes := 0
		if habit, ok := manager.PromptListDetail(promoHabit).Map()[prompt.ID]; ok {
			liked = habit.Liked
			saved = habit.Saved
			likes = int(habit.LikeCount)
		}
		resp.List = append(resp.List, types.Prompt{
			Id:         uint(prompt.ID),
			Title:      prompt.Title,
			Content:    prompt.Content,
			CategoryID: uint(prompt.Category),
			Category:   categoryMap[prompt.Category],
			CreatedBy: types.CreatedBy{
				UserID:   uint(prompt.CreatedBy),
				UserName: userName,
				Avatar:   avatar,
			},
			Likes:     likes,
			Liked:     liked,
			Saved:     saved,
			CreatedAt: prompt.CreatedAt.UTC().Format(utils.TimeLayout),
			UpdatedAt: prompt.UpdatedAt.UTC().Format(utils.TimeLayout),
		})
	}

	return
}

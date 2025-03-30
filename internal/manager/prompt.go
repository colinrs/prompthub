package manager

import (
	"context"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/model"
	"github.com/colinrs/prompthub/pkg/utils"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type PromptDetail struct {
	PromptID      uint      `json:"promptId"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CategoryID    uint      `json:"categoryId"`
	CategoryColor string    `json:"categoryColor"`
	Category      string    `json:"category"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
	CreatedBy     CreatedBy `json:"createdBy"`
}

type CreatedBy struct {
	UserID   uint   `json:"userId"`
	UserName string `json:"userName"`
	Avatar   string `json:"avatar"`
}

type Prompt interface {
	GetPromptDetailByIds(ids []int32) ([]*PromptDetail, error)
	GetPromptCount(ids []int32) (map[int32]*model.PromptsCountTable, error)
	GetPromptLike(ids []int32, userId uint) ([]int32, error)
	GetPromptSave(ids []int32, userId uint) ([]int32, error)
	GetCategoryMap(ids []int32) (map[int32]*model.CategoryTable, error)
}

type promptImpl struct {
	db     *gorm.DB
	svcCtx *svc.ServiceContext
	ctx    context.Context
}

func NewPrompt(ctx context.Context, svcCtx *svc.ServiceContext) Prompt {
	return &promptImpl{
		db:     svcCtx.DB,
		svcCtx: svcCtx,
		ctx:    ctx,
	}
}

func (p *promptImpl) GetPromptDetailByIds(ids []int32) ([]*PromptDetail, error) {
	promptsTable := gen.Use(p.svcCtx.DB).PromptsTable
	promptList, err := promptsTable.Where(promptsTable.PromptsStatus.Eq(config.PromptsStatusNormal), promptsTable.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	categoryIds := make([]int32, 0, len(promptList))
	userIds := make([]int32, 0, len(promptList))
	for _, category := range promptList {
		categoryIds = append(categoryIds, category.Category)
		userIds = append(userIds, category.CreatedBy)
	}
	categoryIds = lo.Uniq(categoryIds)
	categoryMap, err := p.GetCategoryMap(categoryIds)
	if err != nil {
		return nil, err
	}
	userIds = lo.Uniq(userIds)
	userTable := gen.Use(p.svcCtx.DB).UsersTable
	userList, err := userTable.Where(userTable.UserStatus.Eq(config.UserStatusNormal), userTable.ID.In(userIds...)).Find()
	if err != nil {
		return nil, err
	}
	userMap := make(map[int32]*model.UsersTable, len(userList))
	lo.UniqMap(userList, func(user *model.UsersTable, _ int) string {
		userMap[user.ID] = user
		return ""
	})
	list := make([]*PromptDetail, 0, len(promptList))
	for _, prompt := range promptList {
		user, ok := userMap[prompt.CreatedBy]
		userName := ""
		avatar := ""
		if ok {
			userName = user.UserName
			avatar = user.Avatar
		}
		categoryName := ""
		if category, ok := categoryMap[prompt.Category]; ok {
			categoryName = category.CategoryName
		}
		list = append(list, &PromptDetail{
			PromptID:   uint(prompt.ID),
			Title:      prompt.Title,
			Content:    prompt.Content,
			CategoryID: uint(prompt.Category),
			Category:   categoryName,
			CreatedBy: CreatedBy{
				UserID:   uint(prompt.CreatedBy),
				UserName: userName,
				Avatar:   avatar,
			},
			CreatedAt: prompt.CreatedAt.UTC().Format(utils.TimeLayout),
			UpdatedAt: prompt.UpdatedAt.UTC().Format(utils.TimeLayout),
		})
	}
	return list, nil
}

func (p *promptImpl) GetCategoryMap(categoryIds []int32) (map[int32]*model.CategoryTable, error) {
	categoryIds = lo.Uniq(categoryIds)
	categoryTable := gen.Use(p.svcCtx.DB).CategoryTable
	categoryList, err := categoryTable.Where(categoryTable.ID.In(categoryIds...)).Find()
	if err != nil {
		return nil, err
	}
	categoryMap := make(map[int32]*model.CategoryTable, len(categoryList))
	lo.UniqMap(categoryList, func(category *model.CategoryTable, _ int) string {
		categoryMap[category.ID] = category
		return ""
	})
	return categoryMap, nil
}

func (p *promptImpl) GetPromptCount(ids []int32) (map[int32]*model.PromptsCountTable, error) {
	promptsCountTable := gen.Use(p.svcCtx.DB).PromptsCountTable
	promptCountList, err := promptsCountTable.Where(promptsCountTable.PromptsID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	promptCountMap := make(map[int32]*model.PromptsCountTable, len(promptCountList))
	lo.UniqMap(promptCountList, func(prompt *model.PromptsCountTable, _ int) string {
		promptCountMap[prompt.PromptsID] = prompt
		return ""
	})
	return promptCountMap, nil
}

func (p *promptImpl) GetPromptLike(ids []int32, userId uint) ([]int32, error) {
	usersLikeTable := gen.Use(p.svcCtx.DB).UsersLike
	usersLikeList, err := usersLikeTable.Where(usersLikeTable.UserID.Eq(int32(userId)), usersLikeTable.PromptsID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	usersLikeIdList := lo.UniqMap(usersLikeList, func(like *model.UsersLike, _ int) int32 {
		return like.PromptsID
	})
	return usersLikeIdList, nil
}

func (p *promptImpl) GetPromptSave(ids []int32, userId uint) ([]int32, error) {
	usersSaveTable := gen.Use(p.svcCtx.DB).UsersSave
	usersSaveList, err := usersSaveTable.Where(usersSaveTable.UserID.Eq(int32(userId)), usersSaveTable.PromptsID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	usersSaveIdList := lo.UniqMap(usersSaveList, func(like *model.UsersSave, _ int) int32 {
		return like.PromptsID
	})
	return usersSaveIdList, nil
}

package manager

import (
	"context"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/pkg/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type PromptListDetail []*PromptHabit

type PromptHabit struct {
	PromptID    uint  `json:"promptId"`
	Liked       bool  `json:"liked"`
	Saved       bool  `json:"saved"`
	LikeCount   int32 `json:"likeCount"`
	ReviewCount int32 `json:"reviewCount"`
}

func (l PromptListDetail) Len() int {
	return len(l)
}

func (l PromptListDetail) Ids() []int32 {
	return lo.UniqMap(l, func(item *PromptHabit, _ int) int32 {
		return cast.ToInt32(item.PromptID)
	})
}

func (l PromptListDetail) Map() map[int32]*PromptHabit {
	detailMap := make(map[int32]*PromptHabit, len(l))
	lo.UniqMap(l, func(item *PromptHabit, _ int) int32 {
		detailMap[cast.ToInt32(item.PromptID)] = item
		return 0
	})
	return detailMap
}

func (l PromptListDetail) SetSave(ctx context.Context, svcCtx *svc.ServiceContext) {
	p := NewPrompt(ctx, svcCtx)
	userId, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return
	}
	if userId == 0 {
		return
	}
	usersSaveIdList, err := p.GetPromptSave(l.Ids(), userId)
	if err != nil {
		logx.WithContext(ctx).Errorf("GetPromptSave error: %s", err.Error())
		return
	}
	for _, saveId := range usersSaveIdList {
		if lo.Contains(l.Ids(), saveId) {
			l.Map()[saveId].Saved = true
		}
	}
}

func (l PromptListDetail) SetLikes(ctx context.Context, svcCtx *svc.ServiceContext) {
	p := NewPrompt(ctx, svcCtx)
	userId, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return
	}
	if userId == 0 {
		return
	}
	usersLikeIdList, err := p.GetPromptLike(l.Ids(), userId)
	if err != nil {
		logx.WithContext(ctx).Errorf("GetPromptSave error: %s", err.Error())
		return
	}
	for _, likeId := range usersLikeIdList {
		if lo.Contains(l.Ids(), likeId) {
			l.Map()[likeId].Liked = true
		}
	}
}

func (l PromptListDetail) SetCount(ctx context.Context, svcCtx *svc.ServiceContext) {
	p := NewPrompt(ctx, svcCtx)
	promptCountMap, err := p.GetPromptCount(l.Ids())
	if err != nil {
		logx.WithContext(ctx).Errorf("GetPromptSave error: %s", err.Error())
		return
	}
	for _, promptHabit := range l {
		if promptCountMap[int32(promptHabit.PromptID)] != nil {
			promptHabit.LikeCount = promptCountMap[int32(promptHabit.PromptID)].LikeCount
			promptHabit.ReviewCount = promptCountMap[int32(promptHabit.PromptID)].ReviewCount
		}
	}
}

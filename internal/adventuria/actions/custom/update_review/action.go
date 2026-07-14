package update_review

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type reviews interface {
	Save(ctx context.Context, review *model.Review) (*model.Review, error)
	GetByActionID(ctx context.Context, actionId string) (*model.Review, error)
}

var _ model.Action = (*UpdateReview)(nil)

const Type model.ActionType = "update_review"

type UpdateReview struct {
	actions.ActionBase
	reviews reviews
}

func NewDef(reviews reviews) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &UpdateReview{
				ActionBase: actions.NewActionBase(Type),
				reviews:    reviews,
			}
		},
	)
}

func (u *UpdateReview) CanDo(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

type Request struct {
	ActionID string  `json:"action_id" form:"action_id"`
	Comment  *string `json:"comment" form:"comment"`
	Score    *int    `json:"score" form:"score"`
}

func (u *UpdateReview) Do(ctx context.Context, _ *model.Events, _ *model.Player, actionReq model.ActionRequest) (any, error) {
	req, ok := actionReq.(Request)
	if !ok {
		return nil, errors.New("invalid request")
	}

	review, err := u.reviews.GetByActionID(ctx, req.ActionID)
	if err != nil {
		return nil, err
	}

	if req.Comment != nil {
		comment, err := model.NewReviewComment(*req.Comment)
		if err != nil {
			return nil, err
		}
		review.SetComment(comment)
	}

	if req.Score != nil {
		score, err := model.NewReviewScore(*req.Score)
		if err != nil {
			return nil, err
		}
		review.SetScore(score)
	}

	_, err = u.reviews.Save(ctx, review)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

package model

import (
	"errors"

	"github.com/google/uuid"
)

type ReviewData struct {
	Id      string
	Comment ReviewComment
	Score   ReviewScore
}

type Review struct {
	data  ReviewData
	isNew bool
}

func NewReview(id uuid.UUID, comment string, score int) (*Review, error) {
	if id == uuid.Nil {
		return nil, errors.New("review: id cannot be nil")
	}
	reviewComment, err := NewReviewComment(comment)
	if err != nil {
		return nil, err
	}
	reviewScore, err := NewReviewScore(score)
	if err != nil {
		return nil, err
	}

	return &Review{
		data: ReviewData{
			Id:      id.String(),
			Comment: reviewComment,
			Score:   reviewScore,
		},
		isNew: true,
	}, nil
}

func RestoreReview(data ReviewData) *Review {
	return &Review{
		data:  data,
		isNew: false,
	}
}

func (r *Review) IsNew() bool {
	return r.isNew
}

func (r *Review) ID() string {
	return r.data.Id
}

func (r *Review) Comment() ReviewComment {
	return r.data.Comment
}

func (r *Review) SetComment(comment ReviewComment) {
	r.data.Comment = comment
}

func (r *Review) Score() ReviewScore {
	return r.data.Score
}

func (r *Review) SetScore(score ReviewScore) {
	r.data.Score = score
}

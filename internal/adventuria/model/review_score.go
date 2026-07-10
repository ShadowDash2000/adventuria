package model

import (
	"adventuria/internal/adventuria/errs"
)

type ReviewScore uint

func NewReviewScore(score int) (ReviewScore, error) {
	if score < 0 || score > 10 {
		return 0, errs.ErrReviewScoreInvalid
	}

	return ReviewScore(score), nil
}

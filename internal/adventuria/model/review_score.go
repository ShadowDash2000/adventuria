package model

import (
	"adventuria/internal/adventuria/errs"
	"math"
)

type ReviewScore float64

func NewReviewScore(score float64) (ReviewScore, error) {
	if score < 0 || score > 10 {
		return 0, errs.ErrReviewScoreInvalid
	}

	return ReviewScore(math.Trunc(score*10) / 10), nil
}

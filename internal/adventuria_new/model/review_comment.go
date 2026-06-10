package model

import (
	"adventuria/internal/adventuria_new/errs"
)

type ReviewComment string

const MaxCommentLength = 5 * 1024 * 1024 // 5 mb

func NewReviewComment(s string) (ReviewComment, error) {
	if len(s) > MaxCommentLength {
		return "", errs.ErrReviewCommentTooLong
	}

	return ReviewComment(s), nil
}

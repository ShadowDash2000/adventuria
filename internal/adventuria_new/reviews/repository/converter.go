package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/model"

	"github.com/pocketbase/pocketbase/core"
)

func ReviewToRecord(review *model.Review, record *core.Record) {
	record.Id = review.ID()
	record.Set(schema.ReviewSchema.Comment, review.Comment())
	record.Set(schema.ReviewSchema.Score, review.Score())
}

func RecordToReview(record *core.Record) *model.Review {
	return model.RestoreReview(model.ReviewData{
		Id:      record.Id,
		Comment: model.ReviewComment(record.GetString(schema.ReviewSchema.Comment)),
		Score:   model.ReviewScore(record.GetInt(schema.ReviewSchema.Score)),
	})
}

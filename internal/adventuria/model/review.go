package model

type ReviewData struct {
	Id      string
	Comment ReviewComment
	Score   ReviewScore
}

type Review struct {
	data  ReviewData
	isNew bool
}

func NewReview(comment string, score int) (*Review, error) {
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

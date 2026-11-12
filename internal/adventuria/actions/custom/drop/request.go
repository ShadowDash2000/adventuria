package drop

type Request struct {
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
}

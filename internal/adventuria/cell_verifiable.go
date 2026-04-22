package adventuria

type CellVerifiable interface {
	Verify(ctx AppContext, value string) error
}

package adventuria

type EffectVerifiable interface {
	Verify(ctx AppContext, value string) error
}

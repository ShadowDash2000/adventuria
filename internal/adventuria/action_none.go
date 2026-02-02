package adventuria

type NoneAction struct {
	ActionBase
}

func (a *NoneAction) CanDo(ctx ActionContext) bool {
	return ctx.User.LastAction().Type() == ""
}

func (a *NoneAction) Do(_ ActionContext, _ ActionRequest) (*ActionResult, error) {
	return &ActionResult{
		Success: true,
	}, nil
}

func (a *NoneAction) GetVariants(_ ActionContext) any {
	return nil
}

package adventuria

import "adventuria/pkg/result"

type NoneAction struct {
	ActionBase
}

func (a *NoneAction) CanDo(ctx ActionContext) bool {
	return ctx.User.LastAction().Type() == ""
}

func (a *NoneAction) Do(_ ActionContext, _ ActionRequest) (*result.Result, error) {
	return result.Ok(), nil
}

func (a *NoneAction) GetVariants(_ ActionContext) any {
	return nil
}

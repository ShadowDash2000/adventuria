package adventuria

type NoneAction struct {
	ActionBase
}

func (a *NoneAction) CanDo() bool {
	return a.User().LastAction().Type() == ""
}

func (a *NoneAction) Do(_ ActionRequest) (*ActionResult, error) {
	return &ActionResult{
		Success: true,
	}, nil
}

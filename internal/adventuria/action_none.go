package adventuria

type NoneAction struct {
	ActionBase
}

func (a *NoneAction) CanDo(user User) bool {
	return user.LastAction().Type() == ""
}

func (a *NoneAction) Do(_ User, _ ActionRequest) (*ActionResult, error) {
	return &ActionResult{
		Success: true,
	}, nil
}

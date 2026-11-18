package model

type ActionState struct {
	ActivityCompletion *ActionActivityCompletionState
	Activities         *ActionActivitiesState
	UsedItems          ActionUsedItemsState
	ActivityFilter     *ActivityFilter
	Items              *ActionItemsState
}

func (a ActionState) Clone() ActionState {
	return ActionState{
		ActivityCompletion: a.ActivityCompletion.Clone(),
		Activities:         a.Activities.Clone(),
		UsedItems:          a.UsedItems.Clone(),
		ActivityFilter:     a.ActivityFilter.CloneNil(),
		Items:              a.Items.Clone(),
	}
}

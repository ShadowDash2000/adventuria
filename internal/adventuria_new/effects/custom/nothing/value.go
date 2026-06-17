package nothing

const (
	useAfterItemAdd   = "after_item_add"
	useAfterItemUse   = "after_item_use"
	useBeforeGameDone = "before_game_done"
)

var useEvents = map[string]struct{}{
	useAfterItemAdd:   {},
	useAfterItemUse:   {},
	useBeforeGameDone: {},
}

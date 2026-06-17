package go_to_jail

const (
	useAfterItemAdd = "after_item_add"
	useAfterItemUse = "after_item_use"
)

var useEvents = map[string]struct{}{
	useAfterItemAdd: {},
	useAfterItemUse: {},
}

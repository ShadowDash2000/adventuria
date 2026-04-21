package adventuria

type CellRefreshable interface {
	Cell
	RefreshItems(AppContext, Player) error
}

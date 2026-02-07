package adventuria

type CellRefreshable interface {
	Cell
	RefreshItems(AppContext, User) error
}

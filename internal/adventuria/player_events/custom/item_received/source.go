package item_received

type Source string

const (
	SourceItemWheel   Source = "item_wheel"
	SourceShop        Source = "shop"
	SourceAnotherItem Source = "another_item"
	SourceDealer      Source = "dealer"
	SourceDebug       Source = "debug"
)

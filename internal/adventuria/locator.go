package adventuria

import (
	"adventuria/pkg/collections"

	"github.com/pocketbase/pocketbase"
)

type ServiceLocator interface {
	PocketBase() *pocketbase.PocketBase
	Cells() *Cells
	Items() *Items
	Collections() *collections.Collections
	Settings() *Settings
}

type PocketBaseLocator interface {
	PocketBase() *pocketbase.PocketBase
	Collections() *collections.Collections
}

type SettingsLocator interface {
	PocketBase() *pocketbase.PocketBase
	Collections() *collections.Collections
	Settings() *Settings
}

package adventuria

import (
	"adventuria/pkg/cache"
	_ "embed"
	"os"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"

	_ "adventuria/migrations"
)

//go:embed assets/placeholder.jpg
var Placeholder []byte

type GameTest struct {
	BaseGame

	pb *tests.TestApp
}

func NewTestGame() (Game, error) {
	game := &GameTest{
		BaseGame: BaseGame{
			users: cache.NewMemoryCache[string, User](0, true),
		},
	}

	const pbDataDir = "pb_data_test"
	err := os.MkdirAll(pbDataDir, 0777)
	if err != nil {
		return nil, err
	}

	pb, err := tests.NewTestApp(pbDataDir)
	if err != nil {
		return nil, err
	}
	game.pb = pb
	PocketBase = game.pb

	game.Init()

	err = game.createTestUsers()
	if err != nil {
		return nil, err
	}
	err = game.createTestCells()
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (g *GameTest) OnServe(_ func(se *core.ServeEvent) error) {

}

func (g *GameTest) Start() error {
	return nil
}

func (g *GameTest) createTestUsers() error {
	avatar, err := filesystem.NewFileFromBytes(Placeholder, "avatar")
	if err != nil {
		return err
	}

	users := []struct {
		name              string
		password          string
		email             string
		avatar            *filesystem.File
		color             string
		maxInventorySlots int
	}{
		{"user1", "1234567890", "test1@example.com", avatar, "#000000", 3},
		{"user2", "1234567890", "test2@example.com", avatar, "#000000", 3},
	}

	usersCollection, err := GameCollections.Get(CollectionUsers)
	if err != nil {
		return err
	}

	for _, user := range users {
		record := core.NewRecord(usersCollection)
		record.Set("name", user.name)
		record.Set("password", user.password)
		record.Set("email", user.email)
		record.Set("avatar", user.avatar)
		record.Set("color", user.color)
		record.Set("maxInventorySlots", user.maxInventorySlots)
		err = g.pb.Save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GameTest) createTestCells() error {
	cells := []struct {
		cellType string
		name     string
		points   int
		sort     int
	}{
		{"start", "Cell 1 (start)", 10, 100},
		{"game", "Cell 2 (game)", 20, 200},
		{"game", "Cell 3 (game)", 30, 300},
	}

	cellsCollection, err := GameCollections.Get(CollectionCells)
	if err != nil {
		return err
	}

	for _, cell := range cells {
		record := core.NewRecord(cellsCollection)
		record.Set("isActive", true)
		record.Set("type", cell.cellType)
		record.Set("name", cell.name)
		record.Set("points", cell.points)
		record.Set("sort", cell.sort)
		err = g.pb.Save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

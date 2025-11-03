package adventuria

import (
	"adventuria/pkg/cache"
	_ "embed"
	"os"
	"strconv"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/types"

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
	err = game.createTestGames()
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

	for _, user := range users {
		record := core.NewRecord(GameCollections.Get(CollectionUsers))
		record.Set("name", user.name)
		record.Set("password", user.password)
		record.Set("email", user.email)
		record.Set("avatar", user.avatar)
		record.Set("color", user.color)
		record.Set("maxInventorySlots", user.maxInventorySlots)
		record.Set("stats", "{}")
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
		{"shop", "Cell 4 (shop)", 0, 400},
	}

	for _, cell := range cells {
		record := core.NewRecord(GameCollections.Get(CollectionCells))
		record.Set("isActive", true)
		record.Set("type", cell.cellType)
		record.Set("name", cell.name)
		record.Set("points", cell.points)
		record.Set("sort", cell.sort)
		err := g.pb.Save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GameTest) createTestGames() error {
	games := []struct {
		idDb          uint64
		name          string
		releaseDate   types.DateTime
		platforms     []string
		developers    []string
		publishers    []string
		genres        []string
		tags          []string
		steamAppId    uint64
		steamAppPrice int
		campaignTime  float64
	}{
		{idDb: 1, name: "Half-Life", releaseDate: types.NowDateTime(), steamAppId: 1001, steamAppPrice: 1234, campaignTime: 10},
		{idDb: 2, name: "Portal", releaseDate: types.NowDateTime(), steamAppId: 1002, steamAppPrice: 999, campaignTime: 4},
		{idDb: 3, name: "Team Fortress 2", releaseDate: types.NowDateTime(), steamAppId: 1003, steamAppPrice: 0, campaignTime: 0},
		{idDb: 4, name: "Left 4 Dead", releaseDate: types.NowDateTime(), steamAppId: 1004, steamAppPrice: 1499, campaignTime: 6},
		{idDb: 5, name: "Counter-Strike", releaseDate: types.NowDateTime(), steamAppId: 1005, steamAppPrice: 1499, campaignTime: 0},
		{idDb: 6, name: "Dota 2", releaseDate: types.NowDateTime(), steamAppId: 1006, steamAppPrice: 0, campaignTime: 0},
		{idDb: 7, name: "Portal 2", releaseDate: types.NowDateTime(), steamAppId: 1007, steamAppPrice: 1999, campaignTime: 8},
		{idDb: 8, name: "Half-Life 2", releaseDate: types.NowDateTime(), steamAppId: 1008, steamAppPrice: 1999, campaignTime: 15},
		{idDb: 9, name: "Left 4 Dead 2", releaseDate: types.NowDateTime(), steamAppId: 1009, steamAppPrice: 1999, campaignTime: 7},
		{idDb: 10, name: "Counter-Strike 2", releaseDate: types.NowDateTime(), steamAppId: 1010, steamAppPrice: 0, campaignTime: 0},
	}

	for i, game := range games {
		record := core.NewRecord(GameCollections.Get(CollectionGames))
		record.Set("id_db", game.idDb)
		record.Set("name", game.name)
		record.Set("release_date", game.releaseDate)
		record.Set("platforms", game.platforms)
		record.Set("developers", game.developers)
		record.Set("publishers", game.publishers)
		record.Set("genres", game.genres)
		record.Set("tags", game.tags)
		record.Set("steam_app_id", game.steamAppId)
		record.Set("steam_app_price", game.steamAppPrice)
		record.Set("campaign_time", game.campaignTime)
		record.Set("checksum", strconv.Itoa(i))
		err := g.pb.Save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

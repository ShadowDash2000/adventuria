package tests

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/collections"
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
	*adventuria.Game
	ef *adventuria.EffectVerifier
}

func NewGameTest() (*GameTest, error) {
	const pbDataDir = "pb_data_test"
	err := os.MkdirAll(pbDataDir, 0777)
	if err != nil {
		return nil, err
	}

	game := &GameTest{
		Game: adventuria.New(),
	}
	pb, err := tests.NewTestApp(pbDataDir)
	if err != nil {
		return nil, err
	}
	adventuria.PocketBase = pb

	if err = game.init(adventuria.AppContext{App: adventuria.PocketBase}); err != nil {
		return nil, err
	}

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
	err = game.createTestTags()
	if err != nil {
		return nil, err
	}
	err = game.createTestGenres()
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (g *GameTest) init(ctx adventuria.AppContext) error {
	var err error

	adventuria.GameCollections = collections.NewCollections(adventuria.PocketBase)
	adventuria.GameUsers = adventuria.NewUsers(ctx)
	adventuria.GameActions = adventuria.NewActions()
	adventuria.GameCells, err = adventuria.NewCells(ctx)
	if err != nil {
		return err
	}
	adventuria.GameItems, err = adventuria.NewItems(ctx)
	if err != nil {
		return err
	}
	adventuria.GameSettings, err = adventuria.NewSettings(ctx)
	if err != nil {
		return err
	}

	g.ef = adventuria.NewEffectVerifier()
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
		record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionUsers))
		record.Set("name", user.name)
		record.Set("password", user.password)
		record.Set("email", user.email)
		record.Set("avatar", user.avatar)
		record.Set("color", user.color)
		record.Set("maxInventorySlots", user.maxInventorySlots)
		record.Set("stats", "{}")
		err = adventuria.PocketBase.Save(record)
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
		value    string
	}{
		{"start", "Cell 1 (start)", 10, 100, ""},
		{"game", "Cell 2 (game)", 20, 200, ""},
		{"game", "Cell 3 (game)", 30, 300, ""},
		{"shop", "Cell 4 (shop)", 0, 400, "{}"},
	}

	icon, err := filesystem.NewFileFromBytes(Placeholder, "icon")
	if err != nil {
		return err
	}

	for _, cell := range cells {
		record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionCells))
		record.Set("isActive", true)
		record.Set("type", cell.cellType)
		record.Set("name", cell.name)
		record.Set("points", cell.points)
		record.Set("sort", cell.sort)
		record.Set("value", cell.value)
		record.Set("icon", icon)
		err := adventuria.PocketBase.Save(record)
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
		slug          string
		releaseDate   types.DateTime
		platforms     []string
		developers    []string
		publishers    []string
		genres        []string
		tags          []string
		steamAppId    uint64
		steamAppPrice int
	}{
		{idDb: 1, name: "Half-Life", slug: "half-life", releaseDate: types.NowDateTime(), steamAppId: 1001, steamAppPrice: 1234},
		{idDb: 2, name: "Portal", slug: "portal", releaseDate: types.NowDateTime(), steamAppId: 1002, steamAppPrice: 999},
		{idDb: 3, name: "Team Fortress 2", slug: "team-fortress-2", releaseDate: types.NowDateTime(), steamAppId: 1003, steamAppPrice: 0},
		{idDb: 4, name: "Left 4 Dead", slug: "left-4-dead", releaseDate: types.NowDateTime(), steamAppId: 1004, steamAppPrice: 1499},
		{idDb: 5, name: "Counter-Strike", slug: "counter-strike", releaseDate: types.NowDateTime(), steamAppId: 1005, steamAppPrice: 1499},
		{idDb: 6, name: "Dota 2", slug: "dota-2", releaseDate: types.NowDateTime(), steamAppId: 1006, steamAppPrice: 0},
		{idDb: 7, name: "Portal 2", slug: "portal-2", releaseDate: types.NowDateTime(), steamAppId: 1007, steamAppPrice: 1999},
		{idDb: 8, name: "Half-Life 2", slug: "half-life-2", releaseDate: types.NowDateTime(), steamAppId: 1008, steamAppPrice: 1999},
		{idDb: 9, name: "Left 4 Dead 2", slug: "left-4-dead-2", releaseDate: types.NowDateTime(), steamAppId: 1009, steamAppPrice: 1999},
		{idDb: 10, name: "Counter-Strike 2", slug: "counter-strike-2", releaseDate: types.NowDateTime(), steamAppId: 1010, steamAppPrice: 0},
	}

	for i, game := range games {
		record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionActivities))
		record.Set("id_db", game.idDb)
		record.Set("type", adventuria.ActivityTypeGame)
		record.Set("name", game.name)
		record.Set("slug", game.slug)
		record.Set("release_date", game.releaseDate)
		record.Set("platforms", game.platforms)
		record.Set("developers", game.developers)
		record.Set("publishers", game.publishers)
		record.Set("genres", game.genres)
		record.Set("tags", game.tags)
		record.Set("steam_app_id", game.steamAppId)
		record.Set("steam_app_price", game.steamAppPrice)
		record.Set("checksum", strconv.Itoa(i))
		err := adventuria.PocketBase.Save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GameTest) createTestTags() error {
	tags := []struct {
		idDb uint64
		name string
	}{
		{idDb: 1, name: "Adventure"},
	}

	for i, tag := range tags {
		record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionTags))
		record.Set("id_db", tag.idDb)
		record.Set("name", tag.name)
		record.Set("checksum", strconv.Itoa(i))
		err := adventuria.PocketBase.Save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GameTest) createTestGenres() error {
	genres := []struct {
		idDb uint64
		name string
	}{
		{idDb: 1, name: "Action"},
		{idDb: 2, name: "Shooter"},
	}

	for i, genre := range genres {
		record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionGenres))
		record.Set("id_db", genre.idDb)
		record.Set("name", genre.name)
		record.Set("checksum", strconv.Itoa(i))
		err := adventuria.PocketBase.Save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

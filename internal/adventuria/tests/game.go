package tests

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/collections"
	_ "embed"
	"os"
	"strconv"

	"github.com/pocketbase/dbx"
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

	err = game.createTestPlayers()
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
	err = game.createTestSeasons()
	if err != nil {
		return nil, err
	}
	err = game.createTestSettings()
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (g *GameTest) init(ctx adventuria.AppContext) error {
	var err error

	adventuria.GameCollections = collections.NewCollections(adventuria.PocketBase)
	adventuria.GamePlayers = adventuria.NewPlayers(ctx)
	adventuria.GameActions = adventuria.NewActions(ctx)
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

	_ = adventuria.NewInventories(ctx)
	_ = adventuria.NewEffectVerifier(ctx)
	_ = adventuria.NewCellVerifier(ctx)

	adventuria.BindActivitiesHooks(ctx)

	return nil
}

func (g *GameTest) createTestPlayers() error {
	avatar, err := filesystem.NewFileFromBytes(Placeholder, "avatar")
	if err != nil {
		return err
	}

	players := []struct {
		name     string
		password string
		email    string
		avatar   *filesystem.File
		color    string
	}{
		{"player1", "1234567890", "test1@example.com", avatar, "#000000"},
		{"player2", "1234567890", "test2@example.com", avatar, "#000000"},
	}

	for _, player := range players {
		record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionPlayers))
		record.Set(schema.PlayerSchema.Name, player.name)
		record.Set("password", player.password)
		record.Set("email", player.email)
		record.Set(schema.PlayerSchema.Avatar, player.avatar)
		record.Set(schema.PlayerSchema.Color, player.color)
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
		{"start", "Cell 1 (start)", 10, 100, "null"},
		{"game", "Cell 2 (game)", 20, 200, "null"},
		{"game", "Cell 3 (game)", 30, 300, "null"},
		{"shop", "Cell 4 (shop)", 0, 400, "null"},
	}

	icon, err := filesystem.NewFileFromBytes(Placeholder, "icon")
	if err != nil {
		return err
	}

	for _, cell := range cells {
		record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionCells))
		record.Set(schema.CellSchema.Type, cell.cellType)
		record.Set(schema.CellSchema.Name, cell.name)
		record.Set(schema.CellSchema.Points, cell.points)
		record.Set(schema.CellSchema.Sort, cell.sort)
		record.Set(schema.CellSchema.Value, cell.value)
		record.Set(schema.CellSchema.Icon, icon)
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
		record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionActivities))
		record.Set(schema.ActivitySchema.IdDb, game.idDb)
		record.Set(schema.ActivitySchema.Type, adventuria.ActivityTypeGame)
		record.Set(schema.ActivitySchema.Name, game.name)
		record.Set(schema.ActivitySchema.Slug, game.slug)
		record.Set(schema.ActivitySchema.ReleaseDate, game.releaseDate)
		record.Set(schema.ActivitySchema.Platforms, game.platforms)
		record.Set(schema.ActivitySchema.Developers, game.developers)
		record.Set(schema.ActivitySchema.Publishers, game.publishers)
		record.Set(schema.ActivitySchema.Genres, game.genres)
		record.Set(schema.ActivitySchema.Tags, game.tags)
		record.Set(schema.ActivitySchema.SteamAppId, game.steamAppId)
		record.Set(schema.ActivitySchema.SteamAppPrice, game.steamAppPrice)
		record.Set(schema.ActivitySchema.Checksum, strconv.Itoa(i))
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
		record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionTags))
		record.Set(schema.TagSchema.IdDb, tag.idDb)
		record.Set(schema.TagSchema.Name, tag.name)
		record.Set(schema.TagSchema.Checksum, strconv.Itoa(i))
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
		record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionGenres))
		record.Set(schema.GenreSchema.IdDb, genre.idDb)
		record.Set(schema.GenreSchema.Name, genre.name)
		record.Set(schema.GenreSchema.Checksum, strconv.Itoa(i))
		err := adventuria.PocketBase.Save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GameTest) createTestSeasons() error {
	seasons := []struct {
		name            string
		slug            string
		seasonDateStart types.DateTime
		seasonDateEnd   types.DateTime
	}{
		{
			name:            "Season 1",
			slug:            "season-1",
			seasonDateStart: types.NowDateTime(),
			seasonDateEnd:   types.NowDateTime().AddDate(0, 0, 1),
		},
	}

	for _, season := range seasons {
		record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionSeasons))
		record.Set(schema.SeasonSchema.Name, season.name)
		record.Set(schema.SeasonSchema.Slug, season.slug)
		record.Set(schema.SeasonSchema.SeasonDateStart, season.seasonDateStart)
		record.Set(schema.SeasonSchema.SeasonDateEnd, season.seasonDateEnd)
		err := adventuria.PocketBase.Save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GameTest) createTestSettings() error {
	var seasonRecord core.Record
	err := adventuria.PocketBase.RecordQuery(schema.CollectionSeasons).
		Where(dbx.HashExp{schema.SeasonSchema.Slug: "season-1"}).
		One(&seasonRecord)
	if err != nil {
		return err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionSettings))
	record.Set(schema.SettingsSchema.CurrentSeason, seasonRecord.Id)
	record.Set(schema.SettingsSchema.DropsToJail, 2)
	record.Set(schema.SettingsSchema.MaxInventorySlots, 6)
	return adventuria.PocketBase.Save(record)
}

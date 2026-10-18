package v2

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/migrator"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/spf13/cobra"
)

func migrateProgressCommand(
	transactionRunner migrator.TransactionRunner,
	players playersRepository,
	seasons seasonsRepository,
	playerProgress playerProgressService,
	playerStats playerStatsService,
	worlds worldsService,
) *cobra.Command {
	migrator := newProgressMigrator(
		transactionRunner,
		players,
		seasons,
		playerProgress,
		playerStats,
		worlds,
	)

	return &cobra.Command{
		Use:          "progress <base-url>",
		Example:      "migrate-data v2 progress http://127.0.0.1:8080",
		Short:        "Migrates progress from Adventuria v2 through Pocketbase HTTP API",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return migrator.migrate(command.Context(), args[0])
		},
	}
}

const (
	playersPath    = "/api/collections/users/records"
	playersPerPage = 500
)

type progressMigrator struct {
	runInTransaction migrator.TransactionRunner
	players          playersRepository
	seasons          seasonsRepository
	playerProgress   playerProgressService
	playerStats      playerStatsService
	worlds           worldsService
}

func newProgressMigrator(
	transactionRunner migrator.TransactionRunner,
	players playersRepository,
	seasons seasonsRepository,
	playerProgress playerProgressService,
	playerStats playerStatsService,
	worlds worldsService,
) *progressMigrator {
	return &progressMigrator{
		runInTransaction: transactionRunner,
		players:          players,
		seasons:          seasons,
		playerProgress:   playerProgress,
		playerStats:      playerStats,
		worlds:           worlds,
	}
}

func (p *progressMigrator) migrate(ctx context.Context, baseUrl string) error {
	season, err := p.seasons.GetBySlug(ctx, "adventuria-3")
	if err != nil {
		return err
	}

	world, err := p.worlds.GetDefault(ctx)
	if err != nil {
		return err
	}

	res, err := getPlayers(ctx, baseUrl, 1, playersPerPage)
	if err != nil {
		return err
	}

	err = p.createProgresses(ctx, res.Items, season.ID(), world.ID())
	if err != nil {
		return err
	}

	for page := 2; page <= res.TotalPages; page++ {
		res, err := getPlayers(ctx, baseUrl, page, playersPerPage)
		if err != nil {
			return err
		}

		err = p.createProgresses(ctx, res.Items, season.ID(), world.ID())
		if err != nil {
			return err
		}
	}

	return nil
}

type player struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Points      int         `json:"points"`
	CellsPassed int         `json:"cellsPassed"`
	Stats       playerStats `json:"stats"`
	Balance     int         `json:"balance"`
}

type playerStats struct {
	Drops       int `json:"drops"`
	Rerolls     int `json:"rerolls"`
	Finished    int `json:"finished"`
	WasInJail   int `json:"wasInJail"`
	ItemsUsed   int `json:"itemsUsed"`
	DiceRolls   int `json:"diceRolls"`
	MaxDiceRoll int `json:"maxDiceRoll"`
	WheelRolled int `json:"wheelRolled"`
}

type getPlayersQuery struct {
	Page    int    `url:"page"`
	PerPage int    `url:"perPage"`
	Fields  string `url:"fields,omitempty"`
}

type getPlayersResponse struct {
	Items      []player `json:"items"`
	Page       int      `json:"page"`
	PerPage    int      `json:"perPage"`
	TotalItems int      `json:"totalItems"`
	TotalPages int      `json:"totalPages"`
}

func getPlayers(ctx context.Context, baseUrl string, page, perPage int) (*getPlayersResponse, error) {
	q, err := query.Values(getPlayersQuery{
		Page:    page,
		PerPage: perPage,
		Fields:  "id,name,points,cellsPassed,stats,balance",
	})
	if err != nil {
		return nil, fmt.Errorf("encode players query: %w", err)
	}

	requestUrl, err := url.JoinPath(baseUrl, playersPath)
	if err != nil {
		return nil, fmt.Errorf("build players URL: %w", err)
	}

	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return nil, fmt.Errorf("parse players URL: %w", err)
	}
	parsedUrl.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request players: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("request players: unexpected HTTP status %s", res.Status)
	}

	var playersResponse getPlayersResponse
	if err = json.NewDecoder(res.Body).Decode(&playersResponse); err != nil {
		return nil, fmt.Errorf("decode players response: %w", err)
	}

	return &playersResponse, nil
}

func (p *progressMigrator) createProgresses(ctx context.Context, players []player, seasonId, worldId string) error {
	for _, player := range players {
		err := p.runInTransaction(ctx, func(ctx context.Context) error {
			return p.createProgress(ctx, player, seasonId, worldId)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *progressMigrator) createProgress(ctx context.Context, player player, seasonId, worldId string) error {
	playerId, err := p.players.GetIDByName(ctx, player.Name)
	if err != nil {
		if errors.Is(err, errs.ErrPlayerNotFound) {
			return nil
		}
		return err
	}

	progress, err := model.NewPlayerProgress(model.PlayerProgressCreate{
		Player:            playerId,
		Season:            seasonId,
		CurrentWorld:      worldId,
		Energy:            0,
		MaxInventorySlots: 6,
	})
	if err != nil {
		return err
	}

	err = progress.PointsChange(player.Points)
	if err != nil {
		return err
	}
	err = progress.BalanceChange(player.Balance)
	if err != nil {
		return err
	}

	stats, err := model.NewPlayerStats(model.PlayerStatsCreate{
		Player: playerId,
		Season: seasonId,
	})
	if err != nil {
		return err
	}

	stats.CellsPassedChange(player.CellsPassed)
	stats.DropsChange(player.Stats.Drops)
	stats.RerollsChange(player.Stats.Rerolls)
	stats.GamesCompletedChange(player.Stats.Finished)
	stats.WasInJailChange(player.Stats.WasInJail)
	stats.ItemsUsedChange(player.Stats.ItemsUsed)
	stats.DiceRollsChange(player.Stats.DiceRolls)
	stats.MaxDiceRollChange(player.Stats.MaxDiceRoll)
	stats.WheelsRolledChange(player.Stats.WheelRolled)

	_, err = p.playerProgress.Save(ctx, progress)
	if err != nil {
		return err
	}

	_, err = p.playerStats.Save(ctx, stats)
	if err != nil {
		return err
	}

	return nil
}

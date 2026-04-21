package stream_tracker

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"context"
	"errors"
	"os"
	"time"

	streamlive "github.com/ShadowDash2000/is-stream-live"
	"github.com/pocketbase/pocketbase/core"
)

type StreamTracker struct {
	client         *streamlive.StreamLive
	started        bool
	twitchPlayers  map[string]string // twitch_login -> player_id
	youtubePlayers map[string]string // youtube_channel_id -> player_id
}

func NewStreamTracker() (*StreamTracker, error) {
	var clients []streamlive.Client

	twitchClientId, twitchClientIdOk := os.LookupEnv("TWITCH_CLIENT_ID")
	twitchClientSecret, twitchClientSecretOk := os.LookupEnv("TWITCH_CLIENT_SECRET")
	if twitchClientIdOk && twitchClientSecretOk {
		clients = append(clients, streamlive.NewTwitch(twitchClientId, twitchClientSecret))
	}

	youtubeApiKey, youtubeApiKeyOk := os.LookupEnv("YOUTUBE_API_KEY")
	if youtubeApiKeyOk {
		clients = append(clients, streamlive.NewYouTube(youtubeApiKey))
	}

	if len(clients) == 0 {
		return nil, errors.New("stream_tracker: no clients found, expected at least one")
	}

	return &StreamTracker{
		client: streamlive.New(clients...),
	}, nil
}

func (s *StreamTracker) bindHooks(ctx adventuria.AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(schema.CollectionPlayers).BindFunc(func(e *core.RecordEvent) error {
		if client, ok := s.client.GetClient(streamlive.TwitchClientName); ok {
			s.updatePlayerLogin(client, s.twitchPlayers, e.Record.Id, e.Record.GetString("twitch"))
		}
		if client, ok := s.client.GetClient(streamlive.YouTubeClientName); ok {
			s.updatePlayerLogin(client, s.youtubePlayers, e.Record.Id, e.Record.GetString("youtube_channel_id"))
		}
		return e.Next()
	})
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionPlayers).BindFunc(func(e *core.RecordEvent) error {
		if client, ok := s.client.GetClient(streamlive.TwitchClientName); ok {
			s.updatePlayerLogin(client, s.twitchPlayers, e.Record.Id, e.Record.GetString("twitch"))
		}
		if client, ok := s.client.GetClient(streamlive.YouTubeClientName); ok {
			s.updatePlayerLogin(client, s.youtubePlayers, e.Record.Id, e.Record.GetString("youtube_channel_id"))
		}

		return e.Next()
	})
	ctx.App.OnRecordAfterDeleteSuccess(schema.CollectionPlayers).BindFunc(func(e *core.RecordEvent) error {
		twitchLogin := e.Record.GetString("twitch")
		if twitchClient, ok := s.client.GetClient(streamlive.TwitchClientName); ok {
			twitchClient.RemoveLogin(twitchLogin)
		}
		delete(s.twitchPlayers, twitchLogin)

		youtubeChannelId := e.Record.GetString("youtube_channel_id")
		if youtubeClient, ok := s.client.GetClient(streamlive.YouTubeClientName); ok {
			youtubeClient.RemoveLogin(youtubeChannelId)
		}
		delete(s.youtubePlayers, youtubeChannelId)

		return e.Next()
	})
}

func (s *StreamTracker) updatePlayerLogin(client streamlive.Client, players map[string]string, playerId string, newLogin string) {
	if client == nil || players == nil {
		return
	}

	var prevLogin string
	for login, id := range players {
		if id == playerId {
			prevLogin = login
			break
		}
	}

	if prevLogin == newLogin {
		return
	}
	if prevLogin != "" {
		client.RemoveLogin(prevLogin)
		delete(players, prevLogin)
	}
	if newLogin != "" {
		client.AddLogin(newLogin)
		players[newLogin] = playerId
	}
}

func (s *StreamTracker) playerIdForChannel(channel string) (string, bool) {
	if playerId, ok := s.twitchPlayers[channel]; ok {
		return playerId, true
	}
	if playerId, ok := s.youtubePlayers[channel]; ok {
		return playerId, true
	}
	return "", false
}

type playerRecord struct {
	Id               string `db:"id"`
	Twitch           string `db:"twitch"`
	YouTubeChannelId string `db:"youtube_channel_id"`
}

func (s *StreamTracker) fetchPlayers(ctx adventuria.AppContext) ([]playerRecord, error) {
	var players []playerRecord
	err := ctx.App.
		RecordQuery(schema.CollectionPlayers).
		Select(schema.PlayerSchema.Id, schema.PlayerSchema.Twitch, schema.PlayerSchema.YouTubeChannelId).
		All(&players)
	return players, err
}

func (s *StreamTracker) Start(appCtx adventuria.AppContext, ctx context.Context) error {
	if s.started {
		panic("stream_tracker: already started")
	}

	players, err := s.fetchPlayers(appCtx)
	if err != nil {
		return err
	}

	s.twitchPlayers = make(map[string]string, len(players))
	s.youtubePlayers = make(map[string]string, len(players))

	s.bindHooks(appCtx)
	s.client.OnStreamChange(func(e *streamlive.StreamChangeEvent) error {
		return s.onStreamChange(appCtx, e)
	})
	s.client.OnRequestError(func(e *streamlive.RequestErrorEvent) error {
		appCtx.App.Logger().Error("Stream tracker request error", "error", e.Error)
		return e.Next()
	})

	twitchClient, _ := s.client.GetClient(streamlive.TwitchClientName)
	youtubeClient, _ := s.client.GetClient(streamlive.YouTubeClientName)
	for _, player := range players {
		if player.Twitch != "" {
			if twitchClient != nil {
				twitchClient.AddLogin(player.Twitch)
			}
			s.twitchPlayers[player.Twitch] = player.Id
		}

		if player.YouTubeChannelId != "" {
			if youtubeClient != nil {
				youtubeClient.AddLogin(player.YouTubeChannelId)
			}
			s.youtubePlayers[player.YouTubeChannelId] = player.Id
		}
	}

	s.client.StartTracking(ctx, 3*time.Minute)
	s.started = true

	return nil
}

func (s *StreamTracker) onStreamChange(ctx adventuria.AppContext, e *streamlive.StreamChangeEvent) error {
	playerId, ok := s.playerIdForChannel(e.Channel)
	if !ok {
		return e.Next()
	}

	player, err := adventuria.GamePlayers.GetByID(ctx, playerId)
	if err != nil {
		return err
	}

	if player.IsStreamLive() == e.Live {
		return e.Next()
	}

	player.SetIsStreamLive(e.Live)
	err = ctx.App.Save(player.ProxyRecord())
	if err != nil {
		return err
	}

	return e.Next()
}

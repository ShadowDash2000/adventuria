package stream_tracker

import (
	"adventuria/internal/adventuria/model"
	"context"
	"log/slog"
	"os"
	"time"

	streamlive "github.com/ShadowDash2000/is-stream-live"
)

type repository interface {
	UpdateStreamStatusOrSkip(ctx context.Context, playerId string, status bool) (bool, error)
}

type playersRepository interface {
	GetAll(ctx context.Context) ([]*model.PlayerInfo, error)
	NotifyChange(ctx context.Context, id string) error
}

type StreamTracker struct {
	logger            *slog.Logger
	repository        repository
	playersRepository playersRepository
	client            *streamlive.StreamLive
	started           bool
	twitchPlayers     map[string]string // twitch_login -> player_id
	youtubePlayers    map[string]string // youtube_channel_id -> player_id
}

func NewStreamTracker(logger *slog.Logger, repository repository, playersRepository playersRepository) *StreamTracker {
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

	return &StreamTracker{
		logger:            logger,
		repository:        repository,
		playersRepository: playersRepository,
		client:            streamlive.New(clients...),
	}
}

func (s *StreamTracker) Start(ctx context.Context) error {
	if s.started {
		panic("stream_tracker: already started")
	}

	players, err := s.playersRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	s.twitchPlayers = make(map[string]string, len(players))
	s.youtubePlayers = make(map[string]string, len(players))

	s.client.OnStreamChange(func(e *streamlive.StreamChangeEvent) {
		err := s.onStreamChange(e)
		if err != nil {
			s.logger.Error("Stream tracker stream change error", "error", err)
		}
		e.Next()
	})
	s.client.OnRequestError(func(e *streamlive.RequestErrorEvent) {
		s.logger.Error("Stream tracker request error", "errors", e.Errors)
		e.Next()
	})

	twitchClient, _ := s.client.GetClient(streamlive.TwitchClientName)
	youtubeClient, _ := s.client.GetClient(streamlive.YouTubeClientName)
	for _, player := range players {
		if player.Twitch() != "" {
			if twitchClient != nil {
				twitchClient.AddLogin(player.Twitch())
			}
			s.twitchPlayers[player.Twitch()] = player.ID()
		}

		if player.YouTubeChannelId() != "" {
			if youtubeClient != nil {
				youtubeClient.AddLogin(player.YouTubeChannelId())
			}
			s.youtubePlayers[player.YouTubeChannelId()] = player.ID()
		}
	}

	s.client.StartTracking(ctx, 3*time.Minute)
	s.started = true

	return nil
}

func (s *StreamTracker) onStreamChange(e *streamlive.StreamChangeEvent) error {
	playerId, ok := s.playerIdForChannel(e.Channel)
	if !ok {
		return nil
	}

	ok, err := s.repository.UpdateStreamStatusOrSkip(context.Background(), playerId, e.Live)
	if err != nil {
		return err
	}

	if ok {
		err = s.playersRepository.NotifyChange(context.Background(), playerId)
		if err != nil {
			return err
		}
	}

	return nil
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

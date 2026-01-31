package stream_tracker

import (
	"adventuria/internal/adventuria"
	"context"
	"errors"
	"os"
	"time"

	streamlive "github.com/ShadowDash2000/is-stream-live"
	"github.com/pocketbase/pocketbase/core"
)

type StreamTracker struct {
	client       *streamlive.StreamLive
	started      bool
	twitchUsers  map[string]string // twitch_login -> user_id
	youtubeUsers map[string]string // youtube_channel_id -> user_id
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

func (s *StreamTracker) bindHooks() {
	adventuria.PocketBase.OnRecordAfterCreateSuccess(adventuria.CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
		if client, ok := s.client.GetClient(streamlive.TwitchClientName); ok {
			s.updateUserLogin(client, s.twitchUsers, e.Record.Id, e.Record.GetString("twitch"))
		}
		if client, ok := s.client.GetClient(streamlive.YouTubeClientName); ok {
			s.updateUserLogin(client, s.youtubeUsers, e.Record.Id, e.Record.GetString("youtube_channel_id"))
		}
		return e.Next()
	})
	adventuria.PocketBase.OnRecordAfterUpdateSuccess(adventuria.CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
		if client, ok := s.client.GetClient(streamlive.TwitchClientName); ok {
			s.updateUserLogin(client, s.twitchUsers, e.Record.Id, e.Record.GetString("twitch"))
		}
		if client, ok := s.client.GetClient(streamlive.YouTubeClientName); ok {
			s.updateUserLogin(client, s.youtubeUsers, e.Record.Id, e.Record.GetString("youtube_channel_id"))
		}

		return e.Next()
	})
	adventuria.PocketBase.OnRecordAfterDeleteSuccess(adventuria.CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
		twitchLogin := e.Record.GetString("twitch")
		if twitchClient, ok := s.client.GetClient(streamlive.TwitchClientName); ok {
			twitchClient.RemoveLogin(twitchLogin)
		}
		delete(s.twitchUsers, twitchLogin)

		youtubeChannelId := e.Record.GetString("youtube_channel_id")
		if youtubeClient, ok := s.client.GetClient(streamlive.YouTubeClientName); ok {
			youtubeClient.RemoveLogin(youtubeChannelId)
		}
		delete(s.youtubeUsers, youtubeChannelId)

		return e.Next()
	})
}

func (s *StreamTracker) updateUserLogin(client streamlive.Client, users map[string]string, userId string, newLogin string) {
	if client == nil || users == nil {
		return
	}

	var prevLogin string
	for login, id := range users {
		if id == userId {
			prevLogin = login
			break
		}
	}

	if prevLogin == newLogin {
		return
	}
	if prevLogin != "" {
		client.RemoveLogin(prevLogin)
		delete(users, prevLogin)
	}
	if newLogin != "" {
		client.AddLogin(newLogin)
		users[newLogin] = userId
	}
}

func (s *StreamTracker) userIdForChannel(channel string) (string, bool) {
	if userId, ok := s.twitchUsers[channel]; ok {
		return userId, true
	}
	if userId, ok := s.youtubeUsers[channel]; ok {
		return userId, true
	}
	return "", false
}

func (s *StreamTracker) fetchUsers() ([]*core.Record, error) {
	return adventuria.PocketBase.FindAllRecords(adventuria.GameCollections.Get(adventuria.CollectionUsers))
}

func (s *StreamTracker) Start(ctx context.Context) error {
	if s.started {
		panic("stream_tracker: already started")
	}

	users, err := s.fetchUsers()
	if err != nil {
		return err
	}

	s.twitchUsers = make(map[string]string, len(users))
	s.youtubeUsers = make(map[string]string, len(users))

	s.bindHooks()
	s.client.OnStreamChange(s.onStreamChange)
	s.client.OnRequestError(func(e *streamlive.RequestErrorEvent) error {
		adventuria.PocketBase.Logger().Error("Stream tracker request error", "error", e.Error)
		return e.Next()
	})

	twitchClient, _ := s.client.GetClient(streamlive.TwitchClientName)
	youtubeClient, _ := s.client.GetClient(streamlive.YouTubeClientName)
	for _, user := range users {
		twitchLogin := user.GetString("twitch")
		if twitchLogin != "" {
			if twitchClient != nil {
				twitchClient.AddLogin(twitchLogin)
			}
			s.twitchUsers[twitchLogin] = user.Id
		}

		youtubeChannelId := user.GetString("youtube_channel_id")
		if youtubeChannelId != "" {
			if youtubeClient != nil {
				youtubeClient.AddLogin(youtubeChannelId)
			}
			s.youtubeUsers[youtubeChannelId] = user.Id
		}
	}

	s.client.StartTracking(ctx, 3*time.Minute)
	s.started = true

	return nil
}

func (s *StreamTracker) onStreamChange(e *streamlive.StreamChangeEvent) error {
	userId, ok := s.userIdForChannel(e.Channel)
	if !ok {
		return e.Next()
	}

	user, err := adventuria.GameUsers.GetByID(userId)
	if err != nil {
		return err
	}

	user.SetIsStreamLive(e.Live)
	if err = adventuria.PocketBase.Save(user.ProxyRecord()); err != nil {
		return err
	}

	return e.Next()
}

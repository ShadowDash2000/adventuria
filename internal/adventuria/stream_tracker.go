package adventuria

import (
	"context"
	"errors"
	"os"
	"time"

	streamlive "github.com/ShadowDash2000/is-stream-live"
	"github.com/pocketbase/pocketbase/core"
)

type StreamTracker struct {
	client *streamlive.StreamLive
	users  map[string]string // twitch_login -> user_id
}

func NewStreamTracker() (*StreamTracker, error) {
	twitchClientId, ok := os.LookupEnv("TWITCH_CLIENT_ID")
	if !ok {
		return nil, errors.New("stream_tracker: TWITCH_CLIENT_ID not found")
	}
	twitchClientSecret, ok := os.LookupEnv("TWITCH_CLIENT_SECRET")
	if !ok {
		return nil, errors.New("stream_tracker: TWITCH_CLIENT_SECRET not found")
	}

	return &StreamTracker{
		client: streamlive.New(
			streamlive.NewTwitch(twitchClientId, twitchClientSecret),
		),
	}, nil
}

func (s *StreamTracker) bindHooks() {
	PocketBase.OnRecordAfterCreateSuccess(CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
		login := e.Record.GetString("twitch")
		if login != "" {
			s.client.AddLogin(login)
			s.users[login] = e.Record.Id
		}
		return e.Next()
	})
	PocketBase.OnRecordAfterUpdateSuccess(CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
		newLogin := e.Record.GetString("twitch")
		for prevLogin, userId := range s.users {
			if userId != e.Record.Id {
				continue
			}

			if newLogin != prevLogin {
				s.client.RemoveLogin(prevLogin)
				delete(s.users, prevLogin)

				if newLogin != "" {
					s.client.AddLogin(newLogin)
					s.users[newLogin] = userId
				}
			}
			break
		}
		return e.Next()
	})
	PocketBase.OnRecordAfterDeleteSuccess(CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
		login := e.Record.GetString("twitch")
		s.client.RemoveLogin(login)
		delete(s.users, login)
		return e.Next()
	})
}

func (s *StreamTracker) fetchUsers() ([]*core.Record, error) {
	return PocketBase.FindAllRecords(GameCollections.Get(CollectionUsers))
}

func (s *StreamTracker) Start(ctx context.Context) error {
	if s.users != nil {
		panic("stream_tracker: already started")
	}

	users, err := s.fetchUsers()
	if err != nil {
		return err
	}

	s.users = make(map[string]string, len(users))

	s.bindHooks()
	s.client.OnStreamChange(s.onStreamChange)
	s.client.OnRequestError(func(e *streamlive.RequestErrorEvent) error {
		PocketBase.Logger().Error("Stream tracker request error", "error", e.Error)
		return e.Next()
	})

	var logins []string
	for _, user := range users {
		twitchLogin := user.GetString("twitch")
		if twitchLogin == "" {
			continue
		}
		logins = append(logins, twitchLogin)
		s.users[twitchLogin] = user.Id
	}

	s.client.StartTracking(ctx, logins, 3*time.Minute)

	return nil
}

func (s *StreamTracker) onStreamChange(e *streamlive.StreamChangeEvent) error {
	user, err := GameUsers.GetByID(s.users[e.Channel])
	if err != nil {
		return err
	}

	user.SetIsStreamLive(e.Live)
	if err = user.save(); err != nil {
		return err
	}

	return e.Next()
}

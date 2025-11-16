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
	cancel context.CancelFunc
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

	s := &StreamTracker{
		client: streamlive.New(
			streamlive.NewTwitch(twitchClientId, twitchClientSecret),
		),
		users: make(map[string]string),
	}

	s.bindHooks()

	return s, nil
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
		s.client.RemoveLogin(e.Record.GetString("twitch"))
		return e.Next()
	})
}

func (s *StreamTracker) fetchUsers() ([]*core.Record, error) {
	return PocketBase.FindAllRecords(GameCollections.Get(CollectionUsers))
}

func (s *StreamTracker) Start(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}

	users, err := s.fetchUsers()
	if err != nil {
		return err
	}

	s.client.OnStreamChange(s.onStreamChange)

	var logins []string
	for _, user := range users {
		twitchLogin := user.GetString("twitch")
		if twitchLogin == "" {
			continue
		}
		logins = append(logins, twitchLogin)
		s.users[twitchLogin] = user.Id
	}

	if err = s.client.StartTracking(ctx, logins, 5*time.Minute); err != nil {
		s.cancel()
		return err
	}

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

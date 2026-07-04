package stream_tracker

import (
	"adventuria/internal/adventuria/schema"

	streamlive "github.com/ShadowDash2000/is-stream-live"
	"github.com/pocketbase/pocketbase/core"
)

func BindHooks(app core.App, st *StreamTracker) {
	app.OnRecordAfterCreateSuccess(schema.CollectionPlayers).BindFunc(func(e *core.RecordEvent) error {
		if client, ok := st.client.GetClient(streamlive.TwitchClientName); ok {
			st.updatePlayerLogin(client, st.twitchPlayers, e.Record.Id, e.Record.GetString(schema.PlayerSchema.Twitch))
		}
		if client, ok := st.client.GetClient(streamlive.YouTubeClientName); ok {
			st.updatePlayerLogin(client, st.youtubePlayers, e.Record.Id, e.Record.GetString(schema.PlayerSchema.YouTubeChannelId))
		}
		return e.Next()
	})

	app.OnRecordAfterUpdateSuccess(schema.CollectionPlayers).BindFunc(func(e *core.RecordEvent) error {
		if client, ok := st.client.GetClient(streamlive.TwitchClientName); ok {
			st.updatePlayerLogin(client, st.twitchPlayers, e.Record.Id, e.Record.GetString(schema.PlayerSchema.Twitch))
		}
		if client, ok := st.client.GetClient(streamlive.YouTubeClientName); ok {
			st.updatePlayerLogin(client, st.youtubePlayers, e.Record.Id, e.Record.GetString(schema.PlayerSchema.YouTubeChannelId))
		}

		return e.Next()
	})

	app.OnRecordAfterDeleteSuccess(schema.CollectionPlayers).BindFunc(func(e *core.RecordEvent) error {
		twitchLogin := e.Record.GetString(schema.PlayerSchema.Twitch)
		if twitchClient, ok := st.client.GetClient(streamlive.TwitchClientName); ok {
			twitchClient.RemoveLogin(twitchLogin)
		}
		delete(st.twitchPlayers, twitchLogin)

		youtubeChannelId := e.Record.GetString(schema.PlayerSchema.YouTubeChannelId)
		if youtubeClient, ok := st.client.GetClient(streamlive.YouTubeClientName); ok {
			youtubeClient.RemoveLogin(youtubeChannelId)
		}
		delete(st.youtubePlayers, youtubeChannelId)

		return e.Next()
	})
}

package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/model"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToPlayerInfo(record *core.Record) *model.PlayerInfo {
	return model.RestorePlayerInfo(model.PlayerInfoData{
		Id:               record.Id,
		Name:             record.GetString(schema.PlayerSchema.Name),
		Avatar:           record.GetString(schema.PlayerSchema.Avatar),
		Color:            record.GetString(schema.PlayerSchema.Color),
		Twitch:           record.GetString(schema.PlayerSchema.Twitch),
		YouTube:          record.GetString(schema.PlayerSchema.YouTube),
		YouTubeChannelId: record.GetString(schema.PlayerSchema.YouTubeChannelId),
		IsStreamLive:     record.GetBool(schema.PlayerSchema.IsStreamLive),
	})
}

func RecordsToPlayerInfos(records []*core.Record) []*model.PlayerInfo {
	playerInfos := make([]*model.PlayerInfo, len(records))
	for i, record := range records {
		playerInfos[i] = RecordToPlayerInfo(record)
	}
	return playerInfos
}

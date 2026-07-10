package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToSteamSpy(record *core.Record) *model.SteamSpy {
	return model.RestoreSteamSpy(model.SteamSpyData{
		Id:    record.Id,
		IdDb:  record.GetInt(schema.SteamSpySchema.IdDb),
		Name:  record.GetString(schema.SteamSpySchema.Name),
		Price: uint(record.GetInt(schema.SteamSpySchema.Price)),
	})
}

func SteamSpyToRecord(steamSpy *model.SteamSpy, record *core.Record) {
	record.Id = steamSpy.ID()
	record.Set(schema.SteamSpySchema.IdDb, steamSpy.IdDb())
	record.Set(schema.SteamSpySchema.Name, steamSpy.Name())
	record.Set(schema.SteamSpySchema.Price, steamSpy.Price())
}

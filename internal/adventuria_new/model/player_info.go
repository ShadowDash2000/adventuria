package model

type PlayerInfoData struct {
	Id               string
	Name             string
	Avatar           string
	Color            string
	Twitch           string
	YouTube          string
	YouTubeChannelId string
	IsStreamLive     bool
}

type PlayerInfo struct {
	data PlayerInfoData
}

func RestorePlayerInfo(data PlayerInfoData) *PlayerInfo {
	return &PlayerInfo{
		data: data,
	}
}

func (p *PlayerInfo) ID() string {
	return p.data.Id
}

func (p *PlayerInfo) Name() string {
	return p.data.Name
}

func (p *PlayerInfo) Avatar() string {
	return p.data.Avatar
}

func (p *PlayerInfo) Color() string {
	return p.data.Color
}

func (p *PlayerInfo) Twitch() string {
	return p.data.Twitch
}

func (p *PlayerInfo) YouTube() string {
	return p.data.YouTube
}

func (p *PlayerInfo) YouTubeChannelId() string {
	return p.data.YouTubeChannelId
}

func (p *PlayerInfo) IsStreamLive() bool {
	return p.data.IsStreamLive
}

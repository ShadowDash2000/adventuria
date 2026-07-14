package model

type PlayerData struct {
	Id string
}

type Player struct {
	data       PlayerData
	progress   *PlayerProgress
	lastAction *ActionInfo
	stats      *PlayerStats
}

func RestorePlayer(data PlayerData, progress *PlayerProgress, lastAction *ActionInfo, stats *PlayerStats) *Player {
	return &Player{
		data:       data,
		progress:   progress,
		lastAction: lastAction,
		stats:      stats,
	}
}

func (p *Player) ID() string {
	return p.data.Id
}

func (p *Player) Progress() *PlayerProgress {
	return p.progress
}

func (p *Player) SetProgress(progress *PlayerProgress) {
	p.progress = progress
}

func (p *Player) LastAction() *ActionInfo {
	return p.lastAction
}

func (p *Player) SetLastAction(action *ActionInfo) {
	p.lastAction = action
}

func (p *Player) Stats() *PlayerStats {
	return p.stats
}

func (p *Player) SetStats(stats *PlayerStats) {
	p.stats = stats
}

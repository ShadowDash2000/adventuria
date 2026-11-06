package model

type ActivityDoneResultData struct {
	Points        int
	EnergyConsume int
	Coins         int
}

type ActivityDoneResult struct {
	data ActivityDoneResultData
}

func NewActivityDoneResult(data ActivityDoneResultData) *ActivityDoneResult {
	return &ActivityDoneResult{
		data: data,
	}
}

func (d *ActivityDoneResult) Points() int {
	return d.data.Points
}

func (d *ActivityDoneResult) SetPoints(n int) {
	d.data.Points = n
}

func (d *ActivityDoneResult) EnergyConsume() int {
	return d.data.EnergyConsume
}

func (d *ActivityDoneResult) SetEnergyConsume(n int) {
	d.data.EnergyConsume = n
}

func (d *ActivityDoneResult) Coins() int {
	return d.data.Coins
}

func (d *ActivityDoneResult) SetCoins(n int) {
	d.data.Coins = n
}

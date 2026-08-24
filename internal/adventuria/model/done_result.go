package model

type DoneResultData struct {
	Points        int
	EnergyConsume int
	Coins         int
}

type DoneResult struct {
	data DoneResultData
}

func NewDoneResult(data DoneResultData) *DoneResult {
	return &DoneResult{
		data: data,
	}
}

func (d *DoneResult) Points() int {
	return d.data.Points
}

func (d *DoneResult) SetPoints(n int) {
	d.data.Points = n
}

func (d *DoneResult) EnergyConsume() int {
	return d.data.EnergyConsume
}

func (d *DoneResult) SetEnergyConsume(n int) {
	d.data.EnergyConsume = n
}

func (d *DoneResult) Coins() int {
	return d.data.Coins
}

func (d *DoneResult) SetCoins(n int) {
	d.data.Coins = n
}

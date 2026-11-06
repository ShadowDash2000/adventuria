package model

type ActivityDropResultData struct {
	IsBlocked     bool
	IsSafe        bool
	EnergyConsume int
}

type ActivityDropResult struct {
	data ActivityDropResultData
}

func NewActivityDropResult(data ActivityDropResultData) *ActivityDropResult {
	return &ActivityDropResult{
		data: data,
	}
}

func (d *ActivityDropResult) IsBlocked() bool {
	return d.data.IsBlocked
}

func (d *ActivityDropResult) SetIsBlocked(b bool) {
	d.data.IsBlocked = b
}

func (d *ActivityDropResult) IsSafe() bool {
	return d.data.IsSafe
}

func (d *ActivityDropResult) SetIsSafe(b bool) {
	d.data.IsSafe = b
}

func (d *ActivityDropResult) EnergyConsume() int {
	return d.data.EnergyConsume
}

func (d *ActivityDropResult) SetEnergyConsume(n int) {
	d.data.EnergyConsume = n
}

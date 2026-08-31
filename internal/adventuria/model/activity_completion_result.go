package model

type ActivityCompletionResultData struct {
	Points        int
	EnergyConsume int
	Coins         int
}

type ActivityCompletionResult struct {
	data ActivityCompletionResultData
}

func NewActivityCompletionResult(data ActivityCompletionResultData) *ActivityCompletionResult {
	return &ActivityCompletionResult{
		data: data,
	}
}

func (c *ActivityCompletionResult) Points() int {
	return c.data.Points
}

func (c *ActivityCompletionResult) SetPoints(n int) {
	c.data.Points = n
}

func (c *ActivityCompletionResult) EnergyConsume() int {
	return c.data.EnergyConsume
}

func (c *ActivityCompletionResult) SetEnergyConsume(n int) {
	c.data.EnergyConsume = n
}

func (c *ActivityCompletionResult) Coins() int {
	return c.data.Coins
}

func (c *ActivityCompletionResult) SetCoins(n int) {
	c.data.Coins = n
}

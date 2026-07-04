package cheapshark

type CheapSharkResponse struct {
	SteamAppID  uint    `json:"steamAppID"`
	Title       string  `json:"title"`
	NormalPrice float64 `json:"normalPrice"`
}

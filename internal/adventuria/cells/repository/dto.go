package repository

type cellDTO struct {
	Id                       string  `db:"id"`
	Disabled                 bool    `db:"disabled"`
	Sort                     int     `db:"sort"`
	Type                     string  `db:"type"`
	World                    string  `db:"world"`
	Filter                   string  `db:"filter"`
	AudioPreset              string  `db:"audio_preset"`
	Icon                     string  `db:"icon"`
	Name                     string  `db:"name"`
	Points                   int     `db:"points"`
	EnergyConsume            int     `db:"energy_consume"`
	Coins                    int     `db:"coins"`
	Description              string  `db:"description"`
	Color                    string  `db:"color"`
	CantDrop                 bool    `db:"cant_drop"`
	CantReroll               bool    `db:"cant_reroll"`
	IsSafeDrop               bool    `db:"is_safe_drop"`
	IsCustomFilterNotAllowed bool    `db:"is_custom_filter_not_allowed"`
	IsChangeGameNotAllowed   bool    `db:"is_change_game_not_allowed"`
	Value                    *string `db:"value"`

	LocalOrder  int `db:"local_order"`
	GlobalOrder int `db:"global_order"`
}

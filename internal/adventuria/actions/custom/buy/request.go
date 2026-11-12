package buy

import "errors"

type Request struct {
	ItemId string `json:"item_id"`
}

func (r Request) validate() error {
	if r.ItemId == "" {
		return errors.New("buy: item id is required")
	}
	return nil
}
